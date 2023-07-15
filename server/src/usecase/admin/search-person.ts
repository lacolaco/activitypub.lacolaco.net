import { fetchPersonByID } from '@app/activitypub';
import { getTracer } from '@app/tracing';
import { JRDObject } from '@app/webfinger';

const resourceRegexp = /^@(.+)@(.+)$/;

export async function searchPerson(resource: string) {
  return getTracer().startActiveSpan('searchPerson', async (span) => {
    // fetch person id with webfinger
    const [, username, domain] = resource.match(resourceRegexp) ?? [];
    if (domain == null) {
      throw new Error(`invalid resource: ${resource}`);
    }

    const webfingerURL = new URL(`https://${domain}/.well-known/webfinger`);
    const webfingerParams = new URLSearchParams({
      resource: `acct:${username}@${domain}`,
    });
    webfingerURL.search = webfingerParams.toString();

    const res = await fetch(webfingerURL);
    if (!res.ok) {
      throw new Error(`webfinger failed: ${res.status}`);
    }
    const webfinger = (await res.json()) as JRDObject;
    const personURL = webfinger.links.find((link) => link.rel === 'self')?.href;
    if (personURL == null) {
      throw new Error(`person not found: ${resource}`);
    }

    // fetch person with ap
    try {
      return await fetchPersonByID(personURL);
    } catch (e) {
      console.error(e);
      throw new Error(`person not found: ${resource}`);
    }
  });
}
