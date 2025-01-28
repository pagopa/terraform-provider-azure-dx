resource "github_repository" "this" {
  name        = "terraform-provider-azure-dx"
  description = "Terraform provider for DX on Azure, to enhance developer experience"

  visibility = "public"

  allow_auto_merge            = false
  allow_rebase_merge          = true #false
  allow_merge_commit          = true #false
  allow_squash_merge          = true
  squash_merge_commit_title   = "PR_TITLE"
  squash_merge_commit_message = "BLANK"

  delete_branch_on_merge = false

  has_projects = true

  has_issues    = true
  has_downloads = true

  vulnerability_alerts = true
  has_wiki             = true
}
