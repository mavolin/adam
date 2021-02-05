<div align="center">
<h1>adam</h1>

[![Go Reference](https://pkg.go.dev/badge/github.com/mavolin/adam.svg)](https://pkg.go.dev/github.com/mavolin/adam)
![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/mavolin/adam/Test/develop?label=tests)
[![codecov](https://codecov.io/gh/mavolin/adam/branch/develop/graph/badge.svg?token=3qRIAudu4r)](https://codecov.io/gh/mavolin/adam)
[![Go Report Card](https://goreportcard.com/badge/github.com/mavolin/adam)](https://goreportcard.com/report/github.com/mavolin/adam)
[![License](https://img.shields.io/github/license/mavolin/dismock)](https://github.com/mavolin/dismock/blob/v2/LICENSE)
</div>

---

## About

Adam is a bot framework for Discord, built on top of [diamondburned's](https://github.com/diamondburned) library [arikawa](https://github.com/diamondburned).
I originally started working on this because I needed a simple command router with support for localization, but along the way of building it, it turned into a fully-featured bot framework.
You can do everything from a simple `ping` bot to a localized bot with custom commands.

## Main Features

* 🖥️ Typed (variadic) arguments, flags, and out-of-the-box parsing for shellword, and comma-based notations
* 🌍 (optional) support for localization
* 🗒️ Utilities for things like permission handling, emojis, and message and reaction collectors
* ⚡ Error Handling including stack traces
* 👪 Command grouping through modules
* ⏳ Command throttling/cooldowns
* ✏️ Support for message edits
* 🔄 Command overloading through options
* 🤝 Middlewares
* 🛑 Powerful access control system
* 🔌 Custom command sources available at runtime, for things like custom commands
* ✨ Abstracted - Don't like something? Swap it out for a custom implementation

## Getting started

Have a look at the [example bots](./_examples) or use the official [guide](https://go-adam.gitbook.io/adam/) and get your first bot up and running!

## Contributing

Pull requests and issues are appreciated. Check out the [contributing guidelines](./CONTRIBUTING.md) for more information.

## License

Built with ❤️ by [Maximilian von Lindern](https://github.com/mavolin).
Available under the [MIT License](./LICENSE).
