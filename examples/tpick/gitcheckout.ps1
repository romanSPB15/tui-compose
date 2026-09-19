$branch = git branch --format='%(refname:short)' | tpick -t "git checkout"
if ($branch) { git checkout $branch }