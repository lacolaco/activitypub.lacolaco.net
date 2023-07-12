export function getOrigin(host: string) {
  // force https
  return `https://${host}`;
}
