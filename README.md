# DepCharge

DepCharge is a tool designed to help orchestrate commands neccessary for mutli-faceted, larger scale projects.

## Introduction
A medium-to-large sized project (especially when using a microservice architecture) will consist of 3 or more separate repositories, and rely on a variety of package managers depending on the various languages chosen for each service. Typically, these repos must be managed, tracked, and released in some semblance of unison so that dependant service calls can be understood and responded to appropriately.

For large, multi-faceted, development teams, this is a non-issue; a group within the team will manage a service, and assign developers to coordinate. For smaller development teams, this can lead to a handful of devs managing all of the repos. This, in turn, means that a single developer will need to propagating and performing the same commands across all relevant repos. This is a tedious, manual, and error-prone process that occurs every release.

DepCharge is designed to help fix that.

## User Stories

* As a developer, I want to describe all of the dependencies in my project, so that they are all accounted for so that I don't need to manually hunt down all of the repos
    * Describes a config file, dep.yaml
* As a developer, I want to initiate/bootstrap my environment easily, so that I don't need to clone and setup each project/repo
    * Describes the initial `clone` of all, along with potentially a `npm install` and/or `composer install` where neccessary
* As a developer, I want to orchestrate actions across a subset of relevant repos so that a coordinated release can be prepared for multiple services.
    * Describes a method to filter dependencies
* As a developer, I want to define a pattern for subsets of projects to follow so that I can issue commands that act on all at once.
    * This describes the the shortcoming of git submodules, you can't easily perform opperations on just a few.
* As a developer, I want to allow fellow developers to also work within my project subset easily so that we don't need to recreate my setup everytime a new developer steps in.
    * (Trying to describe some form of lock file)
    * Lock files would cause so many merge conflicts if you wanted each branch to have it's own lock that tracked that branch
        * Built-in variables could fix a majority of this
            * Instead of saying something like `- branch: develop`, you could say something like `- branch: $BRANCH` which will look for the active branch, and use that.
                * This implies a parent concept, though. But the dep.yaml has to live somewhere.
