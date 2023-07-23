/**
 * JRD (JSON Resource Descriptor) object
 *
 * @see https://tools.ietf.org/html/rfc7033#section-4.4
 */
export interface JRDObject {
  subject: string;
  links: JRDLinkObject[];
}

/**
 * Webfinger JRD link object
 *
 * @see https://tools.ietf.org/html/rfc7033#section-4.4
 **/
export interface JRDLinkObject {
  rel: string;
  type: string;
  href: string;
}
