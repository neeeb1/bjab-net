---
title: Wing Leader
date: 2026-02-15
slug: wing-leader
tags:
  - web
  - go
description: Community powered, ELO based rating system for bird species.
---
# Wing Leader

Wing Leader is a full-stack web application that allows users to vote on bird species in head-to-head matches and assigns ratings based on the [traditional ELO rating algorithim](https://en.wikipedia.org/wiki/Elo_rating_system).

You can check it out at [wingleader.app](https://wingleader.app)

## Features

- **Type-safe SQL queries utilizing [sqlc](https://github.com/sqlc-dev/sqlc?tab=readme-ov-file)**
- **Database migrations via [goose](https://github.com/pressly/goose)**
- **Structured logging using [zerolog](https://github.com/rs/zerolog)**
- **Lightweight static template generation using [HTMX](https://github.com/bigskysoftware/htmx)**
- **Serverless compute via [Google Cloud Run](https://cloud.google.com/run?hl=en)**
- **Session based vote tracking**
- **IP rate limiting**
- **Concurrent SQL safety**

## Motivation

This project was built as a learning exercise to familiarize myself with Golang backends, PostgreSQL, and CI/CD using Github Actions. More importantly, this project empirically answers an age old question using wisdom of the crowds - which is the best bird?

## Screenshots

![Wing leader homepage](/images/wing-leader-home.png)
  
![Wing leader bird profile](/images/wing-leader-detail.png)