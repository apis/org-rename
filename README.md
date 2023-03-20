# org-rename
Small utility which allows mass update of GitHub organization name for local git repositories

## Command line parameters

```
Usage of org-rename.exe:
  -directory string
        parent directory (default ".")
  -dryRun
        dry run (don't change actual files) (default true)
  -newOrg string
        new organization name (default "PerkinElmerAES")
  -oldOrg string
        old organization name (default "PerkinElmer")
```

## Example run

```
org-rename -directory=. -oldOrg=PerkinElmer -newOrg=PerkinElmerAES -dryRun=false

[INFO] 2023-03-16 17:24:54 - Flag dryRun assigned to 'false'
[INFO] 2023-03-16 17:24:54 - Flag directory assigned to '.'
[INFO] 2023-03-16 17:24:54 - Flag oldOrg assigned to 'PerkinElmer'
[INFO] 2023-03-16 17:24:54 - Flag newOrg assigned to 'PerkinElmerAES'
[INFO] 2023-03-16 17:24:54 - ===============================================
[INFO] 2023-03-16 17:24:54 - Your local GitHub repositories will be updated!
[INFO] 2023-03-16 17:24:54 - ===============================================
[INFO] 2023-03-16 17:24:54 - Analyzing '.git\config'
[INFO] 2023-03-16 17:24:54 - Analyzing 'test\.git\config'
[INFO] 2023-03-16 17:24:54 - Changed '  url = git@github.com:PerkinElmerAES/cds.git'
[INFO] 2023-03-16 17:24:54 - Changed '[lfs "https://github.com/PerkinElmerAES/cds.git/info/lfs"]'
[INFO] 2023-03-16 17:24:54 - Changed '  url = git@github.com:PerkinElmerAES/simba.git'
```
