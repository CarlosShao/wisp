# mock-web fixtures (SPEC-10 §11: mock-web 静态注入载荷页 · 票 22 "compose
# mock-web fixtures: clean page, injection-payload page (for 19/25 tests),
# redirect chain, private-IP host").
#
# These files are served by the mock-web nginx service in compose.test.yml.
# The authoritative acceptance teeth for ticket 22 are the Go suites
# (docs/evidence/s1/22r1-impl.md — httptest-based, CI-portable without
# docker); this service exists for manual/integration red-teaming against a
# realistic HTTP stack, and for tickets 19/25's injection-payload reruns.
#
#   GET http://127.0.0.1:18084/clean.html          -> benign page
#   GET http://127.0.0.1:18084/injection.html      -> D30 attack-script page
#   GET http://127.0.0.1:18084/redirect/1          -> chain 1->2->3->clean
#   GET http://127.0.0.1:18084/redirect/file       -> 302 file:/// (guard must block)
#   GET http://127.0.0.1:18084/redirect/private    -> 302 169.254.169.254 (block)
#   GET http://127.0.0.1:18084/redirect/downgrade  -> 302 http:// (from :18085 TLS origin)
#   GET http://127.0.0.1:18084/huge.bin            -> 4 MiB body (cap + spill)
