export function getGCPProjectID() {
  console.log(JSON.stringify(process.env));
  return process.env.GCP_PROJECT || process.env.GCLOUD_PROJECT || '';
}
