# Contributing to Slacker

First of all, thank you for considering a contribution to Slacker!

On this page you will get an overview of the contribution process for Slacker.
We want to ensure your pull request is merged so when in doubt, please talk to
us!

## Issues

This section addresses what we look for in a good issue, and helps us more
quickly identify and resolve your issue.

### Submitting an Issue

* Please test against the latest commit on the `master` branch. Sometimes
  we've already solved your issue!

* Provide steps and optionally a trivial example we can use to reproduce
  the issue. Please include the actual results you get, and if possible the
  expected results.

* Any examples, errors, or other messages should be provided in text format
  unless it is specifically related to the way Slack is rendering data.

* Remove any sensitive information (tokens, passwords, ip addresses, etc.) from
  your submission.

### Issue Lifecycle

1. Issue is reported

2. A maintainer will triage the issue

3. If it's not critical, the issue may stay inactive for a while until it gets
   picked up. If you feel comfortable trying to address the issue please let us
   know and take a look at the Pull Requests section below.

4. The issue will be resolved via a pull request. We'll reference the issue in
   the pull request to ensure a link exists between the issue and the code that
   fixes it.

5. Issues that we can't reproduce or where the reporter has gone inactive may
   be closed after a period of time at the discretion of the maintainers. If the
   issue is still relevant, we encourage re-opening the issue so it can be
   revisited.

## Pull Requests

This section guides you through making a successful pull request.

### Identifying an Issue

* Pull Requests should address an existing issue. If one doesn't exist, please
  create one. You don't need an issue for trivial PRs like fixing a typo.

* Review the issue and check to see someone else hasn't already begun work on
  this issue. Mention in the comments that you are interested in working on
  this issue.

* A maintainer will reach out to discuss your interest. In some cases there may
  already be some progress towards this issue, or they may have a suggestion on
  how they would like to see this implemented.

* Someone will assign the issue to you to work on. There's no expectation on how
  quickly you will complete this work. We may periodically ask for updates,
  however. Typically we will only do this when there is other interest to
  address the issue. In these situations we do expect you will respond in a
  timely manner.  If you fail to respond after a few requests for an update we
  may re-assign the issue.

### Submitting a PR

* Before submitting your PR for review, run
  [staticcheck](https://staticcheck.io/) against the repository and address any
  identified issues.

* Your pull request should clearly describe what it's accomplishing and how it
  approaches the issue.

* Someone will review your PR and ensure it meets these guidelines. If it does
  not, we will ask you to fix the identified issues.
