#!/bin/bash

PREFIX=Template
AUTOFAC_VERSION=6.0.0
NSUBSTITUTE_VERSION=4.2.2
FLUENT_ASSERTIONS_VERSION=5.10.3
DOTNET_VERSION=net6.0
SOLUTION_NAME=Template-Solution

if [ -z "$1" ]; then
	echo "Give a project name, without Prefix '$PREFIX.'" && exit 1
fi

PROJECT_NAME=${1}
PROJECT_FULL_NAME=${PREFIX}.${PROJECT_NAME}

dotnet new classlib \
	-n "$PROJECT_FULL_NAME" -f "$DOTNET_VERSION" \
	--langVersion latest \
	-o src/"$PROJECT_FULL_NAME"

dotnet add src/"$PROJECT_FULL_NAME" package Autofac -v $AUTOFAC_VERSION

dotnet new xunit \
	-n "$PROJECT_FULL_NAME".Tests \
	-f $DOTNET_VERSION \
	-o src/"$PROJECT_FULL_NAME".Tests

dotnet add src/"$PROJECT_FULL_NAME".Tests package Autofac -v $AUTOFAC_VERSION
dotnet add src/"$PROJECT_FULL_NAME".Tests package NSubstitute -v $NSUBSTITUTE_VERSION
dotnet add src/"$PROJECT_FULL_NAME".Tests package FluentAssertions -v $FLUENT_ASSERTIONS_VERSION

dotnet sln src/"$SOLUTION_NAME".sln add src/"$PROJECT_FULL_NAME"/"$PROJECT_FULL_NAME".csproj
dotnet sln src/"$SOLUTION_NAME".sln add src/"$PROJECT_FULL_NAME".Tests/"$PROJECT_FULL_NAME".Tests.csproj
