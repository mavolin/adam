<div align="center">
<h1>adam</h1>

[![Go Reference](https://pkg.go.dev/badge/github.com/mavolin/adam.svg)](https://pkg.go.dev/github.com/mavolin/adam)
[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/mavolin/adam/Test/develop?label=tests)](https://github.com/mavolin/adam/actions)
[![codecov](https://codecov.io/gh/mavolin/adam/branch/develop/graph/badge.svg?token=3qRIAudu4r)](https://codecov.io/gh/mavolin/adam)
[![Go Report Card](https://goreportcard.com/badge/github.com/mavolin/adam)](https://goreportcard.com/report/github.com/mavolin/adam)
[![License](https://img.shields.io/github/license/mavolin/adam)](https://github.com/mavolin/adam/blob/develop/LICENSE)
</div>

---

## About

Adam is a bot framework for Discord, built on top of [arikawa](https://github.com/diamondburned). I originally started
working on this because I needed a simple command router with support for localization, but along the way of building
it, it turned into a fully-featured bot framework. You can do everything from a simple `ping` bot to a localized bot
with custom commands.

## Main Features

* 🖥️ Typed (variadic) arguments and flags, as well as out-of-the-box parsing for shellword, and delimiter-based notations
* 🌍 Support for localization
* 🗒️ Utilities for permission handling, emojis, and awaiting responses and reactions
* ⚡ Error handling with automatic stack trace generation
* 👪 Command grouping through modules
* ⏳ Command throttling/cooldowns
* ✏️ Support for message edits
* 🤝 Middlewares
* 🛑 Powerful access control system
* 🔌 Custom command sources for commands available at runtime
* ✨ Abstracted - Don't like something? Swap it out for a custom implementation

## Getting Started

Have a look at the [example bots](./_examples) or use the official [guide](https://go-adam.gitbook.io/adam/) and get
your first bot up and running!

## Contributing

Contributions through both pull requests and issues are much appreciated. Check out
the [contributing guidelines](./CONTRIBUTING.md) for more information.

You can also help to localize adam on our [POEditor page](https://poeditor.com/join/project?hash=yLTbnUFjXW).

## License

Built with ❤️ by [Maximilian von Lindern](https://github.com/mavolin). Available under the [MIT License](./LICENSE).
