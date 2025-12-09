# 🧭 AI-Map  
### *The AI-Native Code Navigation Metadata Standard*

**AI-Map** is a lightweight, machine-readable metadata file (`.ai-map.yaml`) that helps AI agents understand and navigate codebases faster and more accurately. It gives LLMs the architectural context they usually lack — without requiring massive embeddings, project-wide scans, or guesswork.

If you’ve ever watched an AI agent wander your repo like a lost intern, this standard is your new best friend.

---

## 🚀 Why AI-Map Exists

Modern AI coding assistants struggle with:

- Multi-repo workspaces  
- Backend + frontend hybrids  
- Complex domains  
- Cross-service interactions  
- Runtime-specific details  
- Critical-path awareness  

AI-Map solves this by giving AI agents a **map**.  
Not documentation.  
Not comments.  
A **machine-readable architectural fingerprint**.

This lets agents:

- Jump directly to the right files  
- Avoid breaking critical paths  
- Produce accurate documentation  
- Understand domain boundaries  
- Generate safer migrations  
- Perform faster refactors  
- Produce fewer hallucinations  

All from a simple YAML file at the repo root.

---

## 📦 What’s Included

AI-Map v1.0 defines a clean, minimal metadata schema that describes:

- **System identity**  
- **Domain + purpose**  
- **Entrypoints**  
- **Data models**  
- **Critical paths**  
- **Internal + external dependencies**  
- **Runtime environment**  
- **Ownership + documentation locations**

This is enough for agents to act meaningfully smarter without burdening developers.

---

## 🗂️ Example `.ai-map.yaml`

```yaml
version: 1

system:
  name: edge-assets
  type: service
  domain: assets
  language: typescript

boundaries:
  entrypoints:
    graphql:
      - src/graphql/resolvers
    http:
      - src/api
  models:
    - src/models
  critical:
    - src/core

dependencies:
  internal:
    - edge-accounts
    - edge-globals
  external:
    - mongodb.atlas
    - redis.cache
    - stripe

ownership:
  team: assets-platform
  slack: "#team-assets"
  docs:
    adr: docs/architecture/adr
    runbook: docs/runbook.md

runtime:
  environment: lambda
  deploys_via: github-actions
  config_paths:
    - infra/config
    - .env.example
```

This is intentionally minimal. Add only what matters.

---

## 📘 Full Specification

The complete AI-Map v1.0 spec is available here:

👉 **[`spec/AI-Map-v1.0.md`](spec/AI-Map-v1.0.md)**

It includes:

- Formal schema  
- Field definitions  
- JSON Schema  
- Tooling guidance  
- Agent routing behavior  
- Extension model  

---

## 🧰 Tooling (Coming Soon)

AI-Map will support:

### **• CLI Tools**
- Validate `.ai-map.yaml`  
- Generate graphs  
- Export metadata for agents  

### **• IDE / Editor Plugins**
- Cursor  
- Neovim  
- VS Code  

### **• MCP Server**
A system-level metadata provider for orchestrating multi-agent workflows.

If you want to help build these, PRs are (or soon will be) welcome.

---

## 🧩 Philosophy

AI-Map follows three core principles:

### **1. Minimality**  
If it doesn’t help an AI agent reason better, it doesn’t belong.

### **2. Stability**  
Specs evolve slowly and intentionally.

### **3. Automation-First**  
Anything tools can derive automatically should be automated, not hand-authored.

---

## 🛠️ Who Is This For?

- Developers using AI agents daily  
- Teams adopting multi-agent systems  
- Projects with complex architectures  
- Multi-repo or monorepo setups  
- Organizations documenting their system boundaries  
- Anyone who wants AI to quit guessing how their code works  

---

## 🤝 Contributing

Contributions are welcome!  
The spec is intentionally small, but tooling, examples, and integrations are all fair game.

Soon you’ll be able to:

- Submit extensions  
- Propose schema evolutions  
- Provide real-world examples  
- Add agent-side integrations  

---

## 📜 License

MIT License — free to use, modify, and integrate into your projects.

---

## 💬 Feedback

Open an issue or discussion — this spec is designed for real-world iteration, and your use cases help shape future versions.

---

If you'd like, I can:

- Add badges (version, spec, status, license)  
- Add a project logo  
- Generate a GitHub org structure  
- Create `.ai-map.yaml` templates for different project types  
- Build the initial example repo layout  
