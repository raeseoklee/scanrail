import { createServer } from "node:http";

const port = Number.parseInt(process.env.PORT || "18080", 10);
const host = process.env.HOST || "127.0.0.1";

const server = createServer((req, res) => {
  if (req.url === "/secure") {
    res.setHeader("Content-Security-Policy", "default-src 'self'");
    res.setHeader("X-Content-Type-Options", "nosniff");
    res.setHeader("X-Frame-Options", "DENY");
    res.setHeader("Referrer-Policy", "no-referrer");
  }

  res.setHeader("Content-Type", "text/plain; charset=utf-8");
  res.end("Scanrail headers demo\n");
});

server.listen(port, host, () => {
  console.log(`Scanrail headers demo listening on http://${host}:${port}`);
});
