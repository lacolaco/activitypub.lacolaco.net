import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../environments/environment';
import { ActivityPubPerson, LocalUser } from './models';

@Injectable({ providedIn: 'root' })
export class AdminApiClient {
  #http = inject(HttpClient);

  async getUsers() {
    return firstValueFrom(this.#http.get<LocalUser[]>(`${environment.backend}/admin/users/list`));
  }

  async getUserByUsername(username: string) {
    return firstValueFrom(this.#http.get<LocalUser>(`${environment.backend}/admin/users/${username}`));
  }

  async searchRemotePerson(resource: string) {
    return firstValueFrom(this.#http.get<ActivityPubPerson>(`${environment.backend}/admin/search/person/${resource}`));
  }

  async getFollowers(username: string) {
    return firstValueFrom(this.#http.get<LocalUser[]>(`${environment.backend}/admin/users/${username}/followers`));
  }
}
