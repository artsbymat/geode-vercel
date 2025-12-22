(function () {
  const theme =
    localStorage.getItem("theme") ||
    (window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light");

  if (theme === "dark") {
    document.documentElement.classList.add("dark");
    const syntaxTheme = document.getElementById("syntax-theme");
    if (syntaxTheme) {
      syntaxTheme.href = "/styles/syntax-dark.css";
    }
  }
})();
