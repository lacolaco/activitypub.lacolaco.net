import { Config } from '@app/domain/config';

export type AppContext = {
  Variables: {
    readonly rsaKeyPair: {
      readonly publicKey: string;
      readonly privateKey: string;
    };
  };
};
