document.addEventListener("DOMContentLoaded", () => {
  const root = document.documentElement;
  const toggle = document.querySelector(".theme-toggle");
  const syntaxThemeLink = document.getElementById("syntax-theme");
  const STORAGE_KEY = "theme";

  function applyTheme(theme) {
    root.classList.toggle("dark", theme === "dark");
    localStorage.setItem(STORAGE_KEY, theme);

    if (syntaxThemeLink) {
      syntaxThemeLink.href =
        theme === "dark"
          ? "/styles/syntax-dark.css"
          : "/styles/syntax-light.css";
    }
  }

  toggle.addEventListener("click", () => {
    const isDark = root.classList.contains("dark");
    applyTheme(isDark ? "light" : "dark");
  });
});
