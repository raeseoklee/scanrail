#!/usr/bin/env node
try {
  require("@scanrail/cli/bin/scanrail.js");
} catch (error) {
  if (error && error.code === "MODULE_NOT_FOUND") {
    console.error("Missing dependency: @scanrail/cli");
    console.error("Try reinstalling with: npm install -g scanrail");
    process.exit(1);
  }
  throw error;
}
