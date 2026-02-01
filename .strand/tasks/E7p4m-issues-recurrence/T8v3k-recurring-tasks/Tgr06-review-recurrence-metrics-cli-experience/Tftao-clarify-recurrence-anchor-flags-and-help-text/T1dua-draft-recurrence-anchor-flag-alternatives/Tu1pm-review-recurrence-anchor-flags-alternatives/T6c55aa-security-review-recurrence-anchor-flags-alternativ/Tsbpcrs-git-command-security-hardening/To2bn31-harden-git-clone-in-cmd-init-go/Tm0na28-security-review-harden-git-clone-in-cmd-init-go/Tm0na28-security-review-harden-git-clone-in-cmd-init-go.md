---
type: review
role: reviewer-security
priority: medium
parent: To2bn31-harden-git-clone-in-cmd-init-go
blockers: []
blocks: []
date_created: 2026-02-01T22:30:23.863087Z
date_edited: 2026-02-01T22:30:23.863087Z
owner_approval: false
completed: false
description: ""
---

# Security review: Harden git clone in cmd/init.go

# Description
Verify that the -- separator correctly prevents flag injection for the preset argument in the git clone command. Confirm that TestInitWithMaliciousPreset adequately covers the threat model.

Delegate concerns to the relevant role via subtasks.
