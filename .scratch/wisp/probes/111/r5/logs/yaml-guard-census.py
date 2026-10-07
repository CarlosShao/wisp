# Ticket 111 r5: YAML parse + per-step guard census for .github/workflows/ci.yml.
# Real instrument, real reading: proves the edited file still PARSES and counts
# the steps that carry an `if:` key. Prints machine-checkable counts only.
# NO go commands, NO CI calls, NO writes.
import sys
import yaml

path = sys.argv[1]
with open(path, encoding="utf-8") as fh:
    doc = yaml.safe_load(fh)

jobs = doc["jobs"]
total_steps = 0
total_guarded = 0
lines = []
co_error = 0
for jname, jbody in jobs.items():
    steps = jbody.get("steps", [])
    guarded = 0
    for st in steps:
        if "if" in st:
            guarded += 1
        if "continue-on-error" in st:
            co_error += 1
    total_steps += len(steps)
    total_guarded += guarded
    lines.append("%s steps=%d guarded=%d runs-on=%s" % (jname, len(steps), guarded, jbody.get("runs-on")))

print("YAML_PARSE=OK")
print("\n".join(lines))
print("TOTAL_STEPS=%d" % total_steps)
print("TOTAL_GUARDED=%d" % total_guarded)
print("CONTINUE_ON_ERROR=%d" % co_error)
for jname, jbody in jobs.items():
    for idx, st in enumerate(jbody.get("steps", []), 1):
        name = st.get("name", "(unnamed: %s)" % st.get("uses", "?"))
        print("%s|%d|if=%s|%s" % (jname, idx, st.get("if", "NONE"), name))
