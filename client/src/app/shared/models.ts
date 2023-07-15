export type LocalUser = {
  id: string;
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
