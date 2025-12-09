# 🧭 AI-Map Specification v1.0  
*A Minimal Metadata Standard for AI-Native Code Navigation*

---

## **Overview**

AI-Map is a lightweight, machine-readable metadata file placed at the root of a repository to help AI agents rapidly understand:

- The purpose and domain of a codebase  
- Architectural boundaries  
- Entrypoints and critical paths  
- Data model locations  
- Internal and external dependencies  
- Runtime and deployment environments  
- Ownership and documentation metadata  

The goal is to allow AI systems to traverse, analyze, refactor, document, and reason about codebases **faster, safer, and with dramatically reduced hallucination**.

AI-Map is designed to be:

- Minimal  
- Stable  
- Easy to maintain  
- Backwards-compatible  
- Extensible  

This is **version 1.0**, intended as a practical baseline for real-world tooling.

---

# 1. File Location

AI-Map metadata MUST reside at the repository root:

```
/.ai-map.yaml
```

If a repository contains multiple independent systems, each should have its own `.ai-map.yaml` file located at the appropriate boundary.

---

# 2. YAML Structure

AI-Map v1.0 defines the following top-level keys:

```yaml
version: 1

system:
  name: string
  type: service | webapp | library | infra | monorepo
  domain: string
  language: string

boundaries:
  entrypoints:
    <protocol>: [paths]
  models: [paths]
  critical: [paths]

dependencies:
  internal: [repo-names]
  external: [services]

ownership:
  team: string
  slack: string
  docs:
    adr: path
    runbook: path

runtime:
  environment: lambda | container | node | browser | worker | cli
  deploys_via: github-actions | cdk | terraform | manual | other
  config_paths: [paths]
```

All fields are optional unless otherwise noted.

---

# 3. Field Definitions

## **3.1 version (required)**

```yaml
version: 1
```

Spec version. Allows future expansion with backward compatibility.

---

## **3.2 system (required)**

Describes the identity of the system.

```yaml
system:
  name: edge-assets
  type: service
  domain: assets
  language: typescript
```

- **name** — canonical system identifier  
- **type** — informs agents how to interpret directory layout  
- **domain** — business or functional domain  
- **language** — primary implementation language  

---

## **3.3 boundaries**

Identifies locations AI should treat as meaningful architectural boundaries.

```yaml
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
```

### **entrypoints**  
Paths initiating system behavior (e.g., API resolvers, controllers, CLI handlers).

### **models**  
Paths defining domain models, schemas, or entity definitions.

### **critical**  
Paths containing essential or high-risk logic that agents should treat with extra caution.

---

## **3.4 dependencies**

Internal and external service dependencies.

```yaml
dependencies:
  internal:
    - edge-accounts
    - edge-globals
  external:
    - mongodb.atlas
    - redis.cache
    - stripe
```

Agents use these to infer cross-system impact and execution context.

---

## **3.5 ownership**

Links system components to human owners and documentation.

```yaml
ownership:
  team: assets-platform
  slack: "#team-assets"
  docs:
    adr: docs/architecture/adr
    runbook: docs/runbook.md
```

This enables documentation generation, review routing, and maintenance tasks.

---

## **3.6 runtime**

Defines execution, configuration, and deployment metadata.

```yaml
runtime:
  environment: lambda
  deploys_via: github-actions
  config_paths:
    - infra/config
    - .env.example
```

Agents use this information to reason about:

- cold starts  
- config changes  
- infra updates  
- deployment side-effects  

---

# 4. Usage: How AI Agents Consume AI-Map

AI-Map does not dictate tooling — instead, it provides a contract for any AI system.

Below are canonical usage patterns.

---

## **4.1 Navigation Agents**

Use AI-Map to:

- Limit search scope to entrypoints, models, and critical paths  
- Resolve dependencies to other repos/services  
- Avoid scanning irrelevant directories  

---

## **4.2 Refactoring Agents**

Use AI-Map to:

- Avoid modifying critical paths unless explicitly allowed  
- Maintain domain logic consistency  
- Understand system type (service vs webapp vs library)  

---

## **4.3 Migration Agents**

Use AI-Map to identify:

- Model definitions  
- Config files  
- Affected internal services  
- Deployment consequences  

---

## **4.4 Documentation Agents**

Can automatically generate:

- README skeletons  
- Architecture diagrams  
- Dependency graphs  
- ADR indexes  
- Onboarding docs  

---

## **4.5 Performance & Security Agents**

Use boundaries + runtime + dependencies for:

- Entrypoint bottleneck detection  
- Security hot spots  
- Dependency-driven threat modeling  

---

# 5. Extensibility

AI-Map encourages extensions using an `extensions` top-level key:

```yaml
extensions:
  ai-flow:
    ignore:
      - src/legacy
    safe_write:
      - src/core
```

Tools may add functionality without breaking compatibility.

---

# 6. Validation Schema

```json
{
  "type": "object",
  "required": ["version", "system"],
  "properties": {
    "version": { "type": "number" },
    "system": {
      "type": "object",
      "required": ["name"],
      "properties": {
        "name": { "type": "string" },
        "type": { "type": "string" },
        "domain": { "type": "string" },
        "language": { "type": "string" }
      }
    },
    "boundaries": {
      "type": "object",
      "properties": {
        "entrypoints": { "type": "object" },
        "models": { "type": "array", "items": { "type": "string" } },
        "critical": { "type": "array", "items": { "type": "string" } }
      }
    },
    "dependencies": {
      "type": "object",
      "properties": {
        "internal": { "type": "array", "items": { "type": "string" } },
        "external": { "type": "array", "items": { "type": "string" } }
      }
    },
    "ownership": {
      "type": "object",
      "properties": {
        "team": { "type": "string" },
        "slack": { "type": "string" },
        "docs": {
          "type": "object",
          "properties": {
            "adr": { "type": "string" },
            "runbook": { "type": "string" }
          }
        }
      }
    },
    "runtime": {
      "type": "object",
      "properties": {
        "environment": { "type": "string" },
        "deploys_via": { "type": "string" },
        "config_paths": { "type": "array", "items": { "type": "string" } }
      }
    }
  }
}
```

---

# 7. Design Principles

AI-Map is guided by three core principles:

### **1. Minimality**  
Only high-signal metadata is included.

### **2. Stability**  
Specs evolve slowly to avoid breaking tooling.

### **3. Automation-First**  
Anything derivable by tools should be automated.

---

# 8. Example File

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

---

# 9. License

MIT License (recommended).
