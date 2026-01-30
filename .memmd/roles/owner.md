# Owner

## Role
Owner (human)

## Responsibilities
- Tasks with this role must never be performed purely by an AI agent. If you are an AI agent, your job is to assist the user by asking them the relevant questions, and helping them to do it, never doing it all yourself.
- Make final decisions about which alternatives or plans to accept.
- Record decisions and final rationale in design docs; ensure alternatives are removed or condensed after the decision and the full document is updated to reflect the owner's decision and any user preferences.
- Approve budgets, timelines, and high-level priorities.
- Resolve cross-track conflicts and unblock teams.

## Process
### When presented with alternatives
- Review Designer-provided alternatives and their pros/cons; state the specific design doc being reviewed or approved.
- Choose one alternative or ask Designer for clarifications.
- Create or accept the Architect's breakdown into epics and approve the plan.

### Interactions with agents
- Agents may create tasks for the Owner and wait for decisions; Owners should acknowledge or delegate.
- When an AI agent unblocks any task, it must start a new session with `mode: "new"` and `async: true` using the exact text `do the next task and you can commit if complete or blocked`.
- When owner tasks are unblocked or pending, AI agents must start a new session with `mode: "new"` and `async: true` using the exact text `do the next --role owner task` so the Owner can resolve all owner tasks in order.

## Notes
- Owner is expected to be a human role — avoid fully automating Owner approvals.
