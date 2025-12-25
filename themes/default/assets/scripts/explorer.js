document.addEventListener("DOMContentLoaded", () => {
  const explorer = document.querySelector(".file-explorer");
  if (!explorer) return;

  // The scrollable container is the parent 'nav' element as per base.css
  const scrollContainer = explorer.closest("nav") || explorer;
  const STORAGE_KEY = "geode:explorer:open";
  const SCROLL_KEY = "geode:explorer:scroll";

  const readOpenKeys = () => {
    try {
      const raw = sessionStorage.getItem(STORAGE_KEY);
      if (!raw) return new Set();
      const arr = JSON.parse(raw);
      if (!Array.isArray(arr)) return new Set();
      return new Set(arr.filter((x) => typeof x === "string"));
    } catch {
      return new Set();
    }
  };

  const writeOpenKeys = (set) => {
    try {
      sessionStorage.setItem(STORAGE_KEY, JSON.stringify(Array.from(set)));
    } catch {
      // ignore
    }
  };

  const openKeys = readOpenKeys();

  // Event Listeners
  explorer.addEventListener("click", (e) => {
    const folder = e.target.closest(".folder");
    if (!folder || !explorer.contains(folder)) return;

    const li = folder.parentElement;
    if (!li) return;

    li.classList.toggle("open");

    const key = li.getAttribute("data-node-key");
    if (key) {
      if (li.classList.contains("open")) {
        openKeys.add(key);
      } else {
        openKeys.delete(key);
      }
      writeOpenKeys(openKeys);
    }
  });

  scrollContainer.addEventListener(
    "scroll",
    () => {
      try {
        sessionStorage.setItem(
          SCROLL_KEY,
          scrollContainer.scrollTop.toString(),
        );
      } catch {
        // ignore
      }
    },
    { passive: true },
  );
});
