# ZeroWall — Developer Guide

This guide is for engineers and contributors looking to extend ZeroWall, integrate with its API, or build custom plugins.

---

## 1. Development Environment Setup
To contribute to the ZeroWall core, you will need:
- A FreeBSD 14 development environment (or a VM).
- **Go 1.22+** for backend daemons (`zwapi`, `zwconfigd`).
- **Node.js 20+** and **npm** for the React web dashboard.
- **Git** for version control.

## 2. Codebase Structure
- `/api`: Go source code for the REST API and management logic.
- `/web-ui`: React + TypeScript source code for the dashboard.
- `/firewall`: Logic for generating `pf.conf` and managing kernel states.
- `/services`: Wrapper scripts and config templates for Unbound, DHCP, etc.
- `/kernel`: Specific FreeBSD kernel patches and `loader.conf` tuning files.
- `/scripts`: Build and deployment automation scripts.

## 3. Extending the API
The ZeroWall API is built using a modular design. To add a new endpoint:
1. Define the resource model in `api/models/`.
2. Add the handler logic in `api/handlers/`.
3. Register the route in `api/router.go`.
4. Run `make generate-docs` to update the Swagger/OpenAPI specification.

## 4. UI Development
The web-ui uses a component-based architecture:
- **Components:** Found in `web-ui/src/components/`. Use the predefined ZeroWall design tokens (Tailwind classes).
- **State:** Managed via Zustand in `web-ui/src/store/`.
- **API Clients:** Auto-generated hooks based on the OpenAPI spec.

To run the dev server with hot-reload:
```sh
cd web-ui
npm install
npm run dev
```

## 5. Plugin System (Beta)
ZeroWall supports a modular plugin system. Plugins can:
- Add pages to the Web UI.
- Register new API endpoints.
- Hook into the `zwconfigd` lifecycle to generate custom service configurations.

See `docs/plugins/` (upcoming) for the full Plugin SDK reference.

## 6. Testing Practices
- **Unit Tests:** All Go logic must have corresponding `_test.go` files.
- **Integration Tests:** Use the scripts in `/tests` to validate firewall rule application and NAT logic.
- **UI Tests:** Playwright is used for end-to-end (E2E) testing of critical dashboard flows.

## 7. Contributing
Please refer to `CONTRIBUTING.md` in the root of the repository for details on our pull request process, coding standards, and code of conduct.

---

*Previous: [Admin Guide](13-configration-guide.md)*
