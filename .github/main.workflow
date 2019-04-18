workflow "greenkeeper assign" {
  on = "push"
  resolves = ["Auto Assign"]
}

action "Filters for GitHub Actions" {
  uses = "actions/bin/filter@master"
  args = "branch supermanner-patch-1"
}

action "Auto Assign" {
  uses = "kentaro-m/auto-assign@master"
  needs = ["Filters for GitHub Actions"]
}
