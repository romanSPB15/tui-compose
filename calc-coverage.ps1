$pkgs = go list ./... | Where-Object { $_ -notmatch 'bench|examples' }
$cover = $pkgs -join ','
go test -coverpkg="$cover" -coverprofile="coverage.txt" $pkgs | Out-Null
go tool cover -func="coverage.txt" | Select-Object -Last 1