import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { RxState, stateful } from '@rx-angular/state';
import { lastValueFrom } from 'rxjs';
import { AppStrokedButton } from '../shared/ui/button';
import { FormFieldModule } from '../shared/ui/form-field';

type ActivityPubPerson = {
  id: string;
  name: string;
  inbox: string;
  summary: string;
  icon?: { url: string };
};

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormFieldModule, AppStrokedButton],
  template: `
    <ng-container *ngIf="state$ | async as state">
      <div class="flex-auto flex flex-col justify-start gap-y-2">
        <h2 class="text-2xl">Search User</h2>
        <form [formGroup]="form" (submit)="searchUser()">
          <app-form-field label="@username@hostname" [showLabel]="true">
            <input type="text" placeholder="@username@hostname" formControlName="userId" />
          </app-form-field>
        </form>

        <div *ngIf="state.searched">
          <div *ngIf="state.person as person" class="flex flex-col items-start rounded-lg bg-panel p-4 shadow">
            <div *ngIf="person.icon">
              <img [src]="person.icon.url" class="w-24 h-24 rounded-lg" />
            </div>
            <div class="flex flex-col items-start py-2">
              <span class="font-bold text-xl">{{ person.name }}</span>
              <span class="text-xs break-all text-gray-600">{{ person.id }}</span>
            </div>
            <div class="py-2" [innerHTML]="person.summary"></div>

            <div class="flex flex-row gap-x-2">
              <button app-stroked-button class="bg-white" (click)="requestFollow(person)">Follow</button>
              <button app-stroked-button class="bg-white" (click)="requestUnfollow(person)">Unfollow</button>
            </div>
          </div>

          <div *ngIf="!state.person">
            <p>Not found</p>
          </div>
        </div>
      </div>
    </ng-container>
  `,
  host: { class: 'flex flex-col gap-y-2 h-full' },
})
export class SearchComponent {
  private readonly http = inject(HttpClient);
  private readonly state = new RxState<{ person: ActivityPubPerson | null; searched: boolean }>();
  readonly state$ = this.state.select().pipe(stateful());

  readonly form = new FormGroup({
    userId: new FormControl('@lacolaco@social.mikutter.hachune.net', { nonNullable: true }),
  });

  ngOnInit() {
    this.state.set({ person: null, searched: false });
  }

  async searchUser() {
    const userId = this.form.getRawValue().userId;
    if (!userId) {
      return;
    }
    try {
      const resp = await lastValueFrom(
        this.http.get<{ user: ActivityPubPerson | null }>(`/api/users/search/${userId}`),
      );
      this.state.set({ person: resp.user, searched: true });
      console.log(resp);
    } catch (e) {
      console.error(e);
    }
  }

  async requestFollow(person: ActivityPubPerson) {
    try {
      await lastValueFrom(this.http.post(`/api/following/create`, { id: person.id }));
    } catch (e) {
      console.error(e);
    }
  }

  async requestUnfollow(person: ActivityPubPerson) {
    try {
      await lastValueFrom(this.http.post(`/api/following/delete`, { id: person.id }));
    } catch (e) {
      console.error(e);
    }
  }
}
