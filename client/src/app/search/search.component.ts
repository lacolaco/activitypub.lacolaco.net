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
          <div *ngIf="state.person as person">
            <p>User Found</p>
            <img *ngIf="person.icon" [src]="person.icon.url" class="w-16 h-16" />
            <p>ID: {{ person.id }}</p>
            <p>Name: {{ person.name }}</p>
            <p>Summary: {{ person.summary }}</p>

            <div class="flex flex-row gap-x-2">
              <button app-stroked-button (click)="requestFollow(person)">Follow</button>
              <button app-stroked-button (click)="requestUnfollow(person)">Unfollow</button>
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
      const person = await lastValueFrom(
        this.http.get<ActivityPubPerson | null>(`/api/users/search`, {
          params: { id: userId },
        }),
      );
      this.state.set({ person, searched: true });
      console.log(person);
    } catch (e) {
      console.error(e);
    }
  }

  async requestFollow(person: ActivityPubPerson) {
    try {
      await lastValueFrom(this.http.post(`/api/users/follow`, { id: person.id }));
    } catch (e) {
      console.error(e);
    }
  }

  async requestUnfollow(person: ActivityPubPerson) {
    try {
      await lastValueFrom(this.http.post(`/api/users/unfollow`, { id: person.id }));
    } catch (e) {
      console.error(e);
    }
  }
}
