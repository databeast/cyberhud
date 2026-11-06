// Progressive enhancement for the mode carousel. Without this script the track
// is still a horizontal scroll container the user can swipe/scroll, <details>
// still expand, and forms still submit — so the page is fully usable with JS
// off. This only adds smooth arrow paging and category jumps.
(function () {
  "use strict";

  function trackFor(id) {
    return document.getElementById(id);
  }

  // Scroll amount per arrow press: roughly the visible width minus one card so
  // there's a little overlap for context.
  function pageDelta(track) {
    var card = track.querySelector(".mode-card");
    var cardW = card ? card.getBoundingClientRect().width + 12 /* gap */ : 240;
    var visible = Math.max(1, Math.floor(track.clientWidth / cardW));
    return cardW * Math.max(1, visible - 1);
  }

  function scrollTrack(id, dir) {
    var track = trackFor(id);
    if (!track) return;
    track.scrollBy({ left: dir * pageDelta(track), behavior: "smooth" });
  }

  // Wire arrows.
  document.querySelectorAll(".car-arrow").forEach(function (btn) {
    var id = btn.getAttribute("data-carousel");
    var dir = btn.classList.contains("car-prev") ? -1 : 1;
    btn.addEventListener("click", function () {
      scrollTrack(id, dir);
    });
  });

  // Wire mode-chip selection: clicking a chip reveals that mode's block in the
  // region's single config card and hides the others. All blocks are already in
  // the DOM, so this is pure show/hide — no server round-trip.
  function selectMode(panelId, mode, chip) {
    var panel = document.getElementById(panelId);
    if (!panel) return;
    panel.querySelectorAll(".mode-config").forEach(function (sec) {
      var match = sec.getAttribute("data-mode") === mode;
      if (match) {
        sec.removeAttribute("hidden");
      } else {
        sec.setAttribute("hidden", "");
      }
    });
    // Reflect selection on the chips within the same track.
    if (chip && chip.parentNode) {
      chip.parentNode.querySelectorAll(".mode-chip").forEach(function (c) {
        c.classList.toggle("selected", c === chip);
      });
    }
  }

  document.querySelectorAll(".mode-chip").forEach(function (chip) {
    // Mark the initially-shown (active) mode as selected.
    if (chip.classList.contains("active")) chip.classList.add("selected");
    chip.addEventListener("click", function () {
      selectMode(chip.getAttribute("data-panel"), chip.getAttribute("data-mode"), chip);
    });
  });

  // Wire category jump chips: scroll the first card of that category into view.
  document.querySelectorAll(".cat-chip").forEach(function (chip) {
    var id = chip.getAttribute("data-carousel");
    var cat = chip.getAttribute("data-cat");
    chip.addEventListener("click", function () {
      var track = trackFor(id);
      if (!track) return;
      var target = track.querySelector('.mode-chip[data-cat="' + cssEscape(cat) + '"]');
      if (!target) return;
      // Scroll the track so the target card aligns to the left edge.
      track.scrollTo({
        left: target.offsetLeft - track.offsetLeft,
        behavior: "smooth",
      });
      // Reflect the active category on the chips.
      chip.parentNode.querySelectorAll(".cat-chip").forEach(function (c) {
        c.classList.toggle("active", c === chip);
      });
    });
  });

  // Minimal CSS.escape fallback for older engines.
  function cssEscape(s) {
    if (window.CSS && CSS.escape) return CSS.escape(s);
    return String(s).replace(/[^a-zA-Z0-9_-]/g, "\\$&");
  }
})();
