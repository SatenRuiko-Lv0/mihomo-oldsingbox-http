# VMess Host normalization test patch

This patch is for protocol-shape testing only. It does not add traffic
camouflage, randomized Host selection, or carrier-billing bypass logic.

The goal is to make VMess HTTP and WebSocket first requests deterministic and
easy to compare in packet captures:

- HTTP and WebSocket both resolve the request Host from the same rule;
- an explicit `Host` header wins over the server address;
- differently-cased Host header keys are removed from normal headers, so only
  one Host line is written;
- HTTP no longer randomly chooses a Host value when multiple values are
  configured;
- WebSocket keeps `Request.Host` and `URL.Host` aligned after a Host override.

For TCP loss and retransmission, the kernel retransmits the same byte stream.
The application-level part that matters for testing is to build the initial
request bytes once, consistently, without duplicate Host fields or per-attempt
Host variation.
