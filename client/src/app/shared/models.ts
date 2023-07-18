export type LocalUser = {
  id: string;
  host: string;
  username: string;
  displayName: string;
  description: string;
  icon: { url: string };
  attachments: Array<{ name: string; value: string }>;
};

export type ActivityPubPerson = {
  id: string;
  name: string;
  inbox: string;
  summary: string;
  icon?: { url: string };
};

export type NewUserParams = {
  host: string;
  username: string;
  displayName: string;
  description: string;
  icon: { url: string };
  attachments: Array<{ name: string; value: string }>;
  url: string;
};
