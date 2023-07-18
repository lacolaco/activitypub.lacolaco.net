import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../environments/environment';
import { ActivityPubPerson, LocalUser, NewUserParams } from './models';

@Injectable({ providedIn: 'root' })
export class AdminApiClient {
  #http = inject(HttpClient);

  async getUsers() {
    return firstValueFrom(this.#http.get<LocalUser[]>(`${environment.backend}/admin/users`));
  }

  async getUserByUsername(hostname: string, username: string) {
    return firstValueFrom(this.#http.get<LocalUser>(`${environment.backend}/admin/users/${hostname}/${username}`));
  }

  async createUser(params: NewUserParams) {
    return firstValueFrom(this.#http.post<LocalUser>(`${environment.backend}/admin/users`, params));
  }

  async searchRemotePerson(resource: string) {
    return firstValueFrom(this.#http.get<ActivityPubPerson>(`${environment.backend}/admin/search/person/${resource}`));
  }

  async getFollowers(username: string) {
    return firstValueFrom(this.#http.get<LocalUser[]>(`${environment.backend}/admin/users/${username}/followers`));
  }

  async postUserNote(user: LocalUser, note: { content: string }) {
    return firstValueFrom(
      this.#http.post<LocalUser>(`${environment.backend}/admin/users/${user.host}/${user.username}/notes`, note),
    );
  }
}
