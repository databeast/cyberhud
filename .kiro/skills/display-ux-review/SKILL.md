---
name: display-ux-review
description: Visually review a cyberhud display mode's rendered snapshot PNG as a code-blind end user, judging whether the render delivers the mode's advertised promise. Use when the user asks to review, critique, or check how a display mode looks, or to verify a layout/visual change to a mode in display/modes/. The image-analysis step costs the user real money and is ONLY run when the user explicitly asks for it.
---

# Display UX Review

Look at what a cyberhud display mode actually renders, and judge it the way the
person holding the device would. This skill exists to beat the "blindness
problem": left to guess, an LLM invents a plausible layout, calls it done, and
the human later finds it looks wrong. Here you render the real pixels and look.

## COST GOVERNANCE — read this first, it overrides everything else

Loading a PNG for visual analysis spends the **user's** real money on every
look. The spend is the user's to control, not yours.

- **Never auto-run the visual look.** No hook, no trigger, no "step 5 happens
  automatically." The costly image-read happens only when the user explicitly
  asks for it this turn.
- **Never offer it as a reflex.** Do not tack "want me to check it?" onto the
  end of every edit. That is banned.
- **You MAY recommend it with judgment.** When you have made significant visual
  changes and a paid look would genuinely be worth it, say so plainly: "this
  changed the layout meaningfully — might be worth a review." Then stop. The
  user decides and the user fires it.
- **Do all the free work first.** Build, render the snapshot, run the free
  pixel-measurement script, read the layout math, converge as far as logic
  takes you — so that the one paid look the user chooses to spend is as likely
  as possible to be right. You make their money go further; you never spend it
  for them.
- **One deliberate look per meaningful unit of work.** Do not glue the look to
  a one-line change. Batch the free convergence; the paid look is a checkpoint,
  not a reflex in the edit loop.
- **Reuse what you have already seen.** Within a session, an image you have
  analyzed stays analyzed — describe it from memory, do not re-load it.

If a step you are about to take would load one or more images, state that it
costs the user money and confirm before doing it — unless the user just asked
for the look this turn.

## The two roles — keep them separate

The honesty of this process depends on splitting the person who wrote the code
from the person who judges the screen.

- **Reviewer (code-blind end user).** Does not read Go, knows nothing about
  layout bridges or row-height metrics. Knows exactly one thing: the mode's
  advertised **name and description** — because that is what made them pick it
  from the menu ("oh, that sounds cool") and turn it on. Their whole question
  is: *is this what I was promised?* They report in plain user voice and never
  propose code fixes.
- **Developer (you, normally).** Receives the reviewer's report as a usability
  bug report, traces root cause, fixes, and updates the mode's spec if the
  intent genuinely changed.

When you run the reviewer pass, **be** the code-blind user. Do not let what you
know about the implementation soften the verdict. If the hostname is cut in
half, the user cannot read it — "the width budget explains it" is not a defense.

## The promise lives in the code: the `# User Experience` spec

Every display mode must declare, in code, what its render is supposed to
deliver — so the review checks against a written contract, not your
assumptions. This spec is a `# User Experience` section appended to the mode's
`display/modes/<mode>/doc.go` (leave the existing architecture docs intact).

If the target mode has no `# User Experience` section, author it first (free
work — no image cost). Template:

```go
// # User Experience
//
// Promise (verbatim from catalog registration in cmd.go / init):
//   Name:        "<mode name>"
//   Description: "<mode summary>"
//
// Class: informational | attract/decorative
//
// Invariants (must hold — a usability defect if violated):
//   - ...
//
// Intentional (do NOT report as defects):
//   - ...
//
// Not verifiable from a single still frame:
//   - ... (animation cadence, blink, scroll timing)
//
// Per-resolution notes (especially extreme aspect ratios / narrow portrait):
//   - ...
```

- **Promise** restates the mode's real catalog name and description *verbatim*.
  The reviewer is handed exactly these lines. If the restatement has drifted
  from the actual `init()` registration, that is itself a finding.
- The buckets map onto the reviewer's discipline: **Invariants** = judge hard;
  **Intentional** = never flag; **Not verifiable** = be honest about still-frame
  limits; **Per-resolution** = where aspect ratios bite.
- **Class** decides the whole lens (see below).

## Judge against the promise, not generic "tidy layout"

Read the Class first. It changes what counts as a defect.

- **Informational** (clock, dashboard, thermal, gauges, wifi, storage, ...):
  legibility, fit, centering, alignment, correct glyphs, no clipping — all
  matter. Irregularity is probably a broken promise.
- **Attract / decorative** (attract_hacking, attract_matrix, attract_bokeh,
  attract_plasma, ticker, ...): **the chaos is the intention.** Deliberate
  disorder, text clipping at the edges, noise, streaming off-screen — that is
  the design, and it is what the user was promised. Do not report it as a bug.
  A single frame cannot show motion, cadence, or timing; judging an animation
  as a fixed layout is a category error.

The classic trap: an attract mode's clipped, streaming text looks exactly like
an informational mode's overflow bug. Same pixels, opposite verdict. The Class
is what tells them apart.

## The loop

1. **Target.** One mode, one style/resolution. The user names it; do not expand
   the scope on your own.
2. **Spec.** Ensure the mode's `doc.go` has a `# User Experience` section;
   author or update it if missing. (Free.)
3. **Edit.** Make the code change — batch related edits rather than one-at-a-time.
4. **Converge for free.** Build. Render just the target snapshot (fast):
   from the mode directory,
   `go test -run 'TestThermalPNGSnapshots/<style>$' -count=1 .`
   (substitute the mode's snapshot test name and the style, e.g. the subtest
   name is the style name like `color-240x240`). Run the free pixel-measurement
   script when geometry is in question (see below). Read the layout math.
5. **Paid look — user-triggered only.** When the user asks, load the PNG at
   `display/modes/<mode>/snapshots/<style>_0001.png` and actually look.
6. **Reviewer pass.** As the code-blind user holding the device: does the render
   deliver the promise? Report observation + how-it-feels + severity.
7. **Developer pass.** Triage the report, fix, update the spec if intent shifted.
8. Iterate 3–7, batching free work between paid looks. When done, regenerate
   the full set with `go test .` (no `-run` filter).

### Rendering gotchas
- The snapshot test does `os.RemoveAll(snapshots/)` first, so a single `-run`
  wipes the other snapshots and regenerates only the target. Harmless: the
  `snapshots/` dirs are gitignored regenerable artifacts, not committed golden
  files. Rebuild everything with `go test .` (no filter).
- Mechanical verification (`testsnapshot.VerifyAll`) only checks the PNG
  decodes, has the right dimensions, and is non-blank. It does **not** judge
  whether the output looks right. That judgment is the whole point of this skill.

## Cost-efficient sampling: one representative per layout code path

Never review "all snapshots" — a single mode can have hundreds of registered
style names (thermal has ~460), and each PNG is a paid image-read. The render
code collapses that huge name space into a small number of structurally
distinct layouts. Review one representative per distinct layout, not per name.

The space is partitioned on two orthogonal axes:

1. **Capability (framework-wide, 6 classes):** mono-slow, mono-fast,
   grayscale-slow, grayscale-fast, color-slow, color-fast (defined in
   `display/style/capability.go`; it is the kebab-case prefix on every style
   name). Capability changes only colorization and refresh/static-ness — **not
   geometry**. So just **two representatives — one color, one mono/grayscale —
   cover both sides of every color-vs-mono rendering branch.** You do not need
   all six.

2. **Aspect band → variant set (a per-mode convention, not a framework
   primitive):** each mode's `display/modes/<mode>/styles/styles_WxH.go` files
   map a resolution's aspect band to a fixed set of variants, and each variant
   points at one shared build function — the build function IS the layout code
   path. For thermal: landscape and square share overview/detail/graph/minimal
   (plus timegraph on fast capabilities); portrait gets thermometer/spark/
   heatmap/leds/avg-thermo/minimal; width<64 or height<64 degenerate to
   minimal-only. Different resolutions within one band are responsive
   re-parameterizations of the same code path, not new layouts.

**The rule:** review one representative per distinct shared build function
(layout code path), times one color + one mono/grayscale panel, across the
aspect bands the mode actually declares — plus one tiny/degenerate panel to
confirm the minimal-only band. For thermal that is roughly ~20 snapshots for
full code-path coverage versus ~460 names. This is the difference between an
affordable review and an unaffordable one.

**This must be recomputed per mode.** The capability axis is shared, but the
variant vocabulary and aspect cutoffs are chosen independently by each mode
(clock, dashboard, gauges each declare their own). There is no central table —
read the target mode's own `styles/` registration (its `styles.go` registry and
`styles_WxH.go` files, and the `adaptiveBuild`/`BuildFn` wiring in
`styles/core.go`) to derive that mode's actual bucket list before sampling.
When you propose a sample to the user, say which code paths each chosen
snapshot exercises, so the coverage (and the cost) is explicit.

## Reviewer output format

Write in the user's voice, but keep it structured and actionable — not cute
prose. For each finding:

- **What I see** — plain description, no code terms.
- **How it feels** — the user's reaction relative to the promise.
- **Severity (user terms):**
  - high — "this isn't what I was promised / I can't use it"
  - medium — "it mostly delivers, but something's off and it bugs me"
  - low — "minor, I'd still enjoy it"

The reviewer **has standing to disagree with the spec.** "The description
promised a readable clock; it technically works but floats weirdly high and
bugs me" is a valid finding — it is measured against the promise, and it also
surfaces when the contract itself is wrong.

## Honesty discipline

- Say what a still frame **cannot** verify (motion, cadence, blink, scroll).
- At small sizes (80px, 240px) individual glyph reads are inference. Separate
  **high-confidence** findings (overflow, wrong text, block position) from
  **low-confidence** ones (exact punctuation character).
- In the developer pass, distinguish symptom from root cause, and do not claim
  a root cause you have not traced to a specific line.
- The user is the authority on intended appearance. The spec codifies it, but
  the user can always correct it.

## Free technique: quantify a visual impression before paying to look

When you suspect a block is mis-centered or clipped, confirm it with numbers
instead of spending an image-read. A tiny throwaway Go program (standard
`image` + `image/png`) that scans rows for the first and last row containing a
non-background pixel (luminance threshold) yields: ink top/bottom rows, ink
height, top vs bottom margin, and ink-center vs panel-center. This turned a
vague "looks high" into "ink center 3.5px above panel center, top margin 3 vs
bottom margin 10" for a clock snapshot. It costs nothing. Use it to converge
geometry first; delete the script afterward (it is a diagnostic, not committed).

## Worked examples (the defect classes this skill exists to catch)

1. **Missing glyph** — thermal `color-240x240`: degree symbol rendered as `?`
   tofu ("52.3?C") and the hottest zone's bar rendered as tofu `[????...]`
   while other bars were fine. Passed every mechanical check; looked broken.
2. **Subtle sub-pixel centering** — clock `grayscale-slow-240x135`: the text
   block rode ~3.5px high. Real but small; caught by the free measurement
   script. Symptom ("assumes a height smaller than the screen") was not the
   root cause (metric-box vs visible-ink centering in shared layout infra).
3. **Gross layout failure + stale content** — dashboard `color-slow-80x160`
   (narrow portrait): text laid out for a wider panel overflowed and clipped
   both edges; a stale "cyberdeck" (old project name) was visible only because
   the render could be read. Narrow/portrait aspect ratios are a weak spot.
4. **The trap** — attract_hacking `mono-slow-264x176`: clipped, streaming text
   that looks like example 3's overflow bug but is the intended Hollywood-hacker
   aesthetic. The chaos is the promise. Correct read: does the effect land, and
   be honest that one frame cannot judge the animation.
