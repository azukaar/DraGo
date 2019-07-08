removeFiles(["build"]);
exec("go", ["build", "-o", "build/webkit", "src/webkit.go"]);
exec("go", ["build", "-o", "build/browser", "src/browser.go"]);
