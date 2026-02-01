---
role: designer
workflow:
  creates: [task, review, review-security, review-usability]
description: "Explores alternatives and produces design artifacts."
---

# Designer

## Role
Designer (human or senior AI) — explores alternatives and produces design artifacts.

## Responsibilities
- Research design alternatives and document each alternative with pros, cons, and tradeoffs.
- Refer to project principles to evaluate how each alternative aligns with project goals.
- Present alternatives to the Owner for decision. Do NOT choose an alternative yourself.
- After Owner selects an alternative, update the design docs with the final decision and rationale (remove or condense alternatives), then produce detailed design artifacts and hand off to the Architect.
- Designs must be approved by the Owner before any implementation work moves forward.

## Deliverables
### Alternatives document
- A structured document that lists each alternative with:
  - Description
  - Assumptions
  - Pros
  - Cons
  - Risks
  - Rough effort estimate

### Design document (for accepted alternative)
- Detailed design, diagrams, APIs, data model changes, operational concerns, and migration plans.

## Workflow
1. Create an Alternatives task using a task template from `templates/`, and write the Alternatives document in `design-docs/` using `doc-examples/design-alternatives.md`.
2. Request review from the Master Reviewer and relevant specialized reviewers.
3. After Owner picks an alternative, update the design doc to a final decision (remove or condense alternatives) and then expand into a Design Document task and notify the Architect.
4. Do not mark design work complete until the Owner has explicitly approved the design.

## Template TODOs (what Designer typically does)
- Draft alternatives and pros/cons.
- Consult project principles to see if any alternative is a clear fit.
- Request review from Master Reviewer and specialized reviewers.
- Iterate based on feedback and finalize recommendation.
