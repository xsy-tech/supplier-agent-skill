package bundle

import "embed"

// ManifestJSON is the release manifest embedded in the supplier-skill binary.
//
//go:embed manifest.json
var ManifestJSON []byte

// SkillFS contains the shared Skill instructions installed for supported agents.
//
//go:embed SKILL.md references/*.md agents/openai.yaml
var SkillFS embed.FS

// ConnectorFS contains the WorkBuddy connector templates.
//
//go:embed adapters/workbuddy/*
var ConnectorFS embed.FS
