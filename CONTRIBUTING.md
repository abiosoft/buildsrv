Contributing 
=============

**[Join us on Slack](https://gophers.slack.com/messages/caddy/)** to chat with other Caddy developers! ([Request an invite](http://bit.ly/go-slack-signup), then join the #caddy channel.)


## The Middleware Registry

The middleware registry is the official list of packages that can be distributed with Caddy on Caddy's website. Registering a package [here](https://github.com/caddyserver/buildsrv/blob/master/features/registry.go) can greatly extend Caddy's functionality and benefit more users.

The word "middleware" technically refers to an HTTP handler that is executed during a request. Caddy directives usually invoke a layer of middleware, but they may also start a background service that does not chain in any HTTP handler. For simplicity, we refer to either kind of package as middleware.


### Writing middleware

This project will accept registrations that meet the requirements listed below. **Please open an issue before starting work on your add-on middleware. State specifically what it will do, where it will go in the middleware chain, and how it will look in the Caddyfile.** We and the community may work with you to fine-tune things. This helps prevent duplicated effort and keeps the project on-task to meet its goals.

Use the [caddydev](https://github.com/caddyserver/caddydev) program to help you test your custom middleware.


### Requirements and Terms

To add your package to the middleware registry, it must meet the following requirements.


- The directive name must one lowercase word, unique and differentiable from other directives. It must be clear and obvious what the directive does, not misleading. This is important to maintain consistency and make it easy for users. Choose carefully, since it cannot be changed except by rare exception.

- Project is and will remain actively maintained and updated. This simply means that any issues or pull requests are responded to within a reasonable amount of time by project owner(s) or collaborator(s), that security fixes are applied as soon as possible, and that the package effectively remains in a stable, working state.

- The package must complement the goals of the Caddy project. Packages that are not in the best interest of the project or its users in general or which carry other implications may not be accepted. Similarly, packages should add functionality best provided by a web server. This can be a gray area, so if you're unsure, open an issue and ask before going to all the trouble.

- Packages that use or rely on free third-party services (which are only free) will be decided on a case-by-case basis.

- Packages that integrate with commercial services (including services with free plans/trials) may be added to the registry for a fee.

- The functionality of a package must be unique among other registered packages and the Caddy core. This is not an app store.

- The package must not use cgo and must be cross-platform.

- Package must be under test using the standard Go `testing` package. Tests should not take more than a few seconds to run.

- The project license must be compatible with the main Caddy project license.


We'll talk you through things if your registration has any problems, so don't worry. But we do reserve the right to reject or delay registrations for any of these or other reasons. Feel free to open an issue if you're not sure if your package would qualify.

Packages that are in the registry may be removed at any time for any reason. Usually the reason will be that the project becomes inactive, irrelevant, or is no longer in-line with the goals of Caddy.


### Registering your middleware

1. Open a pull request to add your middleware's **directive** and **package name** to [the Registry variable in registry.go](https://github.com/caddyserver/buildsrv/blob/master/features/registry.go).

2. That's it! We'll take a look.

By submitting a pull request, you verify that you are the project owner or that the project owner has given permission to integrate the package with Caddy. You also understand and agree to comply with the terms and requirements. Adding packages to the registry does not imply endorsement by Caddy, but you may say that your package is available as Caddy middleware.
