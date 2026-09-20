#!/usr/bin/env python3
"""Throwaway generator for ticket-11 golden fixtures. Only the .sse output is
committed. Every dialect carries the SAME semantics as the canonical chat
fixtures so the shared adaptertest table asserts identical expectations."""
import json
import os

DIR = os.path.join("internal", "llm", "testdata", "golden")
C = lambda o: json.dumps(o, separators=(",", ":"))


def dl(payloads):
    return ["data: " + (p if isinstance(p, str) else C(p)) for p in payloads]


def el(pairs):
    out = []
    for name, payload in pairs:
        out.append("event: " + name)
        out.append("data: " + (payload if isinstance(payload, str) else C(payload)))
        out.append("")
    return out


def write(name, lines, scenario):
    all_lines = ["# wisp golden sse v1", "# @scenario " + scenario] + lines
    # An event-style fixture must END with a blank line: that blank line is
    # the SSE dispatcher for the final event (message_stop / [DONE]).
    if all_lines[-1] == "":
        body = "\n".join(all_lines) + "\n"
    else:
        body = "\n".join(all_lines).rstrip("\n") + "\n"
    with open(os.path.join(DIR, name + ".sse"), "w", encoding="utf-8", newline="\n") as f:
        f.write(body)
    print("wrote ", name)


# ------------------------------------------------- chat (new scenarios only)
def c_chunk(delta, finish=None):
    return {"choices": [{"index": 0, "delta": delta, "finish_reason": finish}]}


def c_usage(i, o, cached=None):
    u = {"prompt_tokens": i, "completion_tokens": o, "total_tokens": i + o}
    if cached is not None:
        u["prompt_tokens_details"] = {"cached_tokens": cached}
    return {"choices": [], "usage": u}


write("chat-disconnect-tools", ["# @response 200"] + dl([
    c_chunk({"role": "assistant", "content": ""}),
    c_chunk({"content": "before the tool call"}),
    c_chunk({"tool_calls": [{"index": 0, "id": "call_x1", "type": "function",
                             "function": {"name": "get_weather",
                                          "arguments": "{\"city\":"}}]}),
    c_chunk({"tool_calls": [{"index": 0, "function": {"arguments": "Zh"}}]}),
]), "mid-stream disconnect with an unclosed tool call")

write("chat-thinking", ["# @response 200"] + dl([
    c_chunk({"role": "assistant", "content": "", "reasoning_content": "Counting the dots."}),
    c_chunk({"content": "four"}),
    c_chunk({}, "stop"),
    c_usage(12, 8),
    "[DONE]",
]), "reasoning_content deltas then text")

write("chat-mid-stream-error", ["# @response 200"] + dl([
    c_chunk({"role": "assistant", "content": ""}),
    c_chunk({"content": "up to the error"}),
    {"error": {"message": "backend exploded mid-stream", "type": "server_error",
               "code": "internal_error"}},
]), "error payload inside a 200 stream")

# ----------------------------------------------------------- anthropic side
def a_msg_start(usage):
    msg = {"id": "msg_mock1", "type": "message", "role": "assistant",
           "model": "claude-mock", "content": []}
    if usage is not None:
        msg["usage"] = usage
    return ("message_start", {"type": "message_start", "message": msg})


def a_block_start(index, block):
    return ("content_block_start", {"type": "content_block_start", "index": index,
                                    "content_block": block})


def a_delta(index, delta):
    return ("content_block_delta", {"type": "content_block_delta", "index": index,
                                    "delta": delta})


def a_block_stop(index):
    return ("content_block_stop", {"type": "content_block_stop", "index": index})


def a_msg_delta(stop, usage):
    o = {"type": "message_delta",
         "delta": {"stop_reason": stop, "stop_sequence": None}}
    if usage is not None:
        o["usage"] = usage
    return ("message_delta", o)


A_STOP = ("message_stop", {"type": "message_stop"})


def a_text(index, text):
    return [a_block_start(index, {"type": "text", "text": ""}),
            a_delta(index, {"type": "text_delta", "text": text}),
            a_block_stop(index)]


def a_tool(index, cid, name, args_json):
    # same split point as the canonical chat fixture (first 10 bytes), and an
    # empty object stays ONE fragment, so every dialect streams the identical
    # argument fragments.
    if args_json == "{}":
        return [a_block_start(index, {"type": "tool_use", "id": cid, "name": name, "input": {}}),
                a_delta(index, {"type": "input_json_delta", "partial_json": args_json}),
                a_block_stop(index)]
    half = 10
    return [a_block_start(index, {"type": "tool_use", "id": cid, "name": name, "input": {}}),
            a_delta(index, {"type": "input_json_delta", "partial_json": args_json[:half]}),
            a_delta(index, {"type": "input_json_delta", "partial_json": args_json[half:]}),
            a_block_stop(index)]


write("anthropic-tool-call", ["# @response 200"] + el(
    [a_msg_start({"input_tokens": 21, "cache_read_input_tokens": 8,
                  "cache_creation_input_tokens": 0, "output_tokens": 1})]
    + a_text(0, "Let me check.")
    + a_tool(1, "call_a1", "get_weather", '{"city":"Zhuhai"}')
    + a_tool(2, "call_b2", "get_time", "{}")
    + [a_msg_delta("end_turn", {"output_tokens": 9}), A_STOP]),
    "two fragmented tool calls + text + usage with cache reads")
# stop reason must be tool_use for the tool-call scenario
with open(os.path.join(DIR, "anthropic-tool-call.sse"), encoding="utf-8") as f:
    txt = f.read()
with open(os.path.join(DIR, "anthropic-tool-call.sse"), "w", encoding="utf-8", newline="\n") as f:
    f.write(txt.replace('"stop_reason":"end_turn"', '"stop_reason":"tool_use"'))

write("anthropic-disconnect", ["# @response 200"] + el(
    [a_msg_start({"input_tokens": 10, "output_tokens": 1}),
     a_block_start(0, {"type": "text", "text": ""}),
     a_delta(0, {"type": "text_delta", "text": "partial text kept"})]),
    "stream ends after partial text, no message_stop")

write("anthropic-disconnect-tools", ["# @response 200"] + el(
    [a_msg_start({"input_tokens": 10, "output_tokens": 1})]
    + a_text(0, "before the tool call")
    + [a_block_start(1, {"type": "tool_use", "id": "call_x1", "name": "get_weather", "input": {}}),
       a_delta(1, {"type": "input_json_delta", "partial_json": '{"city":"Zh'})]),
    "stream ends with an open tool_use block")

write("anthropic-usage-multichunk", ["# @response 200"] + el(
    [a_msg_start({"input_tokens": 10, "output_tokens": 4})]
    + a_text(0, "Hello there")
    + [a_msg_delta("end_turn", {"output_tokens": 4}),
       ("message_delta", {"type": "message_delta", "delta": {},
                          "usage": {"output_tokens": 7}}),
       A_STOP]),
    "usage reported twice aggregates by max-merge")

write("anthropic-thinking", ["# @response 200"] + el(
    [a_msg_start({"input_tokens": 12, "output_tokens": 1}),
     a_block_start(0, {"type": "thinking", "thinking": ""}),
     a_delta(0, {"type": "thinking_delta", "thinking": "Counting the dots."}),
     a_delta(0, {"type": "signature_delta", "signature": "AbC="}),
     a_block_stop(0)]
    + a_text(1, "four")
    + [a_msg_delta("end_turn", {"output_tokens": 8}), A_STOP]),
    "extended thinking surfaces as ReasoningDelta")

write("anthropic-max-tokens", ["# @response 200"] + el(
    [a_msg_start({"input_tokens": 10, "output_tokens": 1})]
    + a_text(0, "This answer is cut mid-sentence because the token budget ran")
    + [a_msg_delta("max_tokens", {"output_tokens": 6}), A_STOP]),
    "max_tokens truncation")

write("anthropic-missing-usage", ["# @response 200"] + el(
    [a_msg_start(None)] + a_text(0, "hi") + [a_msg_delta("end_turn", None), A_STOP]),
    "complete stream with no usage block")

write("anthropic-no-terminal", ["# @response 200"] + el(
    [a_msg_start({"input_tokens": 4, "output_tokens": 1})] + a_text(0, "hi") + [A_STOP]),
    "message_stop present but no message_delta (no stop_reason)")

write("anthropic-mid-stream-error", ["# @response 200"] + el(
    [a_msg_start({"input_tokens": 4, "output_tokens": 1})]
    + a_text(0, "up to the error")
    + [("error", {"type": "error", "error": {"type": "overloaded_error",
                                            "message": "Overloaded"}})]),
    "error event inside a 200 stream")

write("anthropic-long-text", ["# @latency 40", "# @response 200"] + el(
    [a_msg_start({"input_tokens": 6, "output_tokens": 1}),
     a_block_start(0, {"type": "text", "text": ""})]
    + [a_delta(0, {"type": "text_delta", "text": "word%d " % i}) for i in range(30)]
    + [a_msg_delta("end_turn", {"output_tokens": 30}), A_STOP]),
    "30 paced chunks for the cancellation case")

# ----------------------------------------------------------- responses side
def r_created():
    return ("response.created", {"type": "response.created", "response": {
        "id": "resp_mock1", "object": "response", "status": "in_progress",
        "model": "gpt-mock"}})


def r_item_added(index, item):
    return ("response.output_item.added", {"type": "response.output_item.added",
                                          "output_index": index, "item": item})


def r_item_done(index, item):
    return ("response.output_item.done", {"type": "response.output_item.done",
                                         "output_index": index, "item": item})


def r_text_delta(item_id, text):
    return ("response.output_text.delta", {"type": "response.output_text.delta",
                                           "item_id": item_id, "delta": text})


def r_args_delta(item_id, part):
    return ("response.function_call_arguments.delta",
            {"type": "response.function_call_arguments.delta", "item_id": item_id,
             "delta": part})


def r_usage(i, o, cached=None):
    u = {"input_tokens": i, "output_tokens": o}
    if cached is not None:
        u["input_tokens_details"] = {"cached_tokens": cached}
    return u


def r_completed(usage, status="completed", extra=None):
    resp = {"id": "resp_mock1", "object": "response", "status": status}
    if usage is not None:
        resp["usage"] = usage
    if extra:
        resp.update(extra)
    return ("response.completed", {"type": "response.completed", "response": resp})


MSG_ITEM = lambda iid: {"id": iid, "type": "message", "role": "assistant",
                        "status": "in_progress", "content": []}
FC_ITEM = lambda iid, cid, name, args="": {"id": iid, "type": "function_call",
                                            "call_id": cid, "name": name,
                                            "arguments": args}
DONE = ("DONE", "[DONE]")

write("responses-tool-call", ["# @response 200"] + el(
    [r_created(),
     r_item_added(0, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "Let me check."),
     r_item_done(0, {"id": "msg_out1", "type": "message", "role": "assistant",
                     "status": "completed",
                     "content": [{"type": "output_text", "text": "Let me check."}]}),
     r_item_added(1, FC_ITEM("fc_1", "call_a1", "get_weather")),
     r_args_delta("fc_1", '{"city":"Zh'),
     r_args_delta("fc_1", 'uhai"}'),
     r_item_done(1, FC_ITEM("fc_1", "call_a1", "get_weather", '{"city":"Zhuhai"}')),
     r_item_added(2, FC_ITEM("fc_2", "call_b2", "get_time")),
     r_args_delta("fc_2", "{}"),
     r_item_done(2, FC_ITEM("fc_2", "call_b2", "get_time", "{}")),
     r_completed(r_usage(21, 9, 8)),
     DONE]),
    "two fragmented tool calls + text + usage with cache reads")

write("responses-disconnect", ["# @response 200"] + el(
    [r_created(), r_item_added(0, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "partial text kept")]),
    "stream ends after partial text, no terminal and no [DONE]")

write("responses-disconnect-tools", ["# @response 200"] + el(
    [r_created(), r_item_added(0, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "before the tool call"),
     r_item_added(1, FC_ITEM("fc_1", "call_x1", "get_weather")),
     r_args_delta("fc_1", '{"city":"Zh')]),
    "stream ends with an open function_call item")

write("responses-usage-multichunk", ["# @response 200"] + el(
    [r_created(), r_item_added(0, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "Hello there"),
     r_completed(r_usage(10, 4)),
     ("response.usage", {"type": "response.usage",
                         "response": {"id": "resp_mock1", "usage": r_usage(10, 7)}}),
     DONE]),
    "usage reported twice (completed + extra usage frame) aggregates by max-merge")

write("responses-thinking", ["# @response 200"] + el(
    [r_created(),
     r_item_added(0, {"id": "rs_1", "type": "reasoning", "summary": [], "content": []}),
     ("response.reasoning_summary_text.delta",
      {"type": "response.reasoning_summary_text.delta", "item_id": "rs_1",
       "output_index": 0, "summary_index": 0, "delta": "Counting the dots."}),
     r_item_done(0, {"id": "rs_1", "type": "reasoning",
                     "summary": [{"type": "summary_text", "text": "Counting the dots."}]}),
     r_item_added(1, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "four"),
     r_completed(r_usage(12, 8)),
     DONE]),
    "reasoning summary deltas surface as ReasoningDelta")

write("responses-max-tokens", ["# @response 200"] + el(
    [r_created(), r_item_added(0, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "This answer is cut mid-sentence because the token budget ran"),
     r_completed(r_usage(10, 6), status="incomplete",
                 extra={"incomplete_details": {"reason": "max_output_tokens"}}),
     DONE]),
    "incomplete response with max_output_tokens reason")

write("responses-missing-usage", ["# @response 200"] + el(
    [r_created(), r_item_added(0, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "hi"),
     ("response.completed", {"type": "response.completed",
                             "response": {"id": "resp_mock1", "status": "completed"}}),
     DONE]),
    "completed response carrying no usage object")

write("responses-no-terminal", ["# @response 200"] + el(
    [r_created(), r_item_added(0, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "hi"),
     ("response.output_text.done", {"type": "response.output_text.done",
                                    "item_id": "msg_out1", "text": "hi"}),
     DONE]),
    "[DONE] arrives without response.completed (no status)")

write("responses-mid-stream-error", ["# @response 200"] + el(
    [r_created(), r_item_added(0, MSG_ITEM("msg_out1")),
     r_text_delta("msg_out1", "up to the error"),
     ("response.failed", {"type": "response.failed", "response": {
         "id": "resp_mock1", "status": "failed",
         "error": {"type": "server_error", "message": "backend exploded mid-stream"}}})]),
    "response.failed inside a 200 stream")

write("responses-long-text", ["# @latency 40", "# @response 200"] + el(
    [r_created(), r_item_added(0, MSG_ITEM("msg_out1"))]
    + [("response.output_text.delta", {"type": "response.output_text.delta",
                                       "item_id": "msg_out1", "delta": "word%d " % i})
       for i in range(30)]
    + [r_completed(r_usage(6, 30)), DONE]),
    "30 paced chunks for the cancellation case")

# ---------------------------------------------------- ladder fixtures (new)
def fault_section(status, msg, etype, code, retry_after=None):
    hdr = "# @response %d" % status
    if retry_after:
        hdr += " retry-after:" + retry_after
    return [hdr, C({"error": {"message": msg, "type": etype, "code": code}})]


RATE = fault_section(429, "Rate limited, retry later", "rate_limit_error",
                     "rate_limit_exceeded", "1")
OK_TAIL_A = el([a_msg_start({"input_tokens": 5, "output_tokens": 1})]
               + a_text(0, "ok after retries")
               + [a_msg_delta("end_turn", {"output_tokens": 3}), A_STOP])
OK_TAIL_R = el([r_created(), r_item_added(0, MSG_ITEM("msg_out1")),
                r_text_delta("msg_out1", "ok after retries"),
                r_completed(r_usage(5, 3)), DONE])

write("anthropic-backoff-429", RATE + RATE + ["# @response 200"] + OK_TAIL_A,
      "two 429s honoring retry-after=1s then success")
write("responses-backoff-429", RATE + RATE + ["# @response 200"] + OK_TAIL_R,
      "two 429s honoring retry-after=1s then success")

E500 = fault_section(500, "upstream exploded", "server_error", "internal_error")
write("anthropic-provider-500", E500 * 5, "five 500s: initial + 3 retries all fail")
write("responses-provider-500", E500 * 5, "five 500s: initial + 3 retries all fail")

A401 = fault_section(401, "Invalid API key", "authentication_error", "invalid_api_key")
write("anthropic-unauthorized", A401, "401: no retry, Unconfigured semantics")
write("responses-unauthorized", A401, "401: no retry, Unconfigured semantics")

print("all fixtures written")
