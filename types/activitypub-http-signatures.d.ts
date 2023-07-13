/**
 * Type definitions for activitypub-http-signatures
 */
declare module 'activitypub-http-signatures' {
  export class Sha256Signer {
    constructor(params: { publicKeyId: string; privateKey: string; headerNames?: string[] });

    sign(params: { url: string; method: string; headers: Record<string, string> }): string;
  }

  export class Parser {
    parse(params: { headers: Record<string, string>; method: string; url: string }): Sha256Signature;
  }

  export class Sha256Signature {
    readonly keyId: string;

    verify(key: string): boolean;
  }

  const parser: Parser;
  export default parser;
}
