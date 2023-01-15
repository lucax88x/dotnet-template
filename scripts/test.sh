#!/bin/bash
#
dotnet test src/Template-Solution.sln \
	--configuration Debug \
	--filter LocalOnly!=true \
	--logger trx \
	--logger "console;verbosity=quiet" \
	--verbosity normal \
	--no-build --no-restore
