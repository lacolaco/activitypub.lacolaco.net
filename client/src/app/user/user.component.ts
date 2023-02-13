import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { map, switchMap } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { RxState, stateful } from '@rx-angular/state';

export type LocalUser = {
  id: string;
  name: string;
  description: string;
  icon: { url: string };
};

@Component({
  selector: 'app-user',
  standalone: true,
  imports: [CommonModule],
  template: `
    <ng-container *ngIf="state$ | async as state">
      <div *ngIf="state.user as user" class="flex flex-col items-start rounded-lg bg-panel p-4 shadow">
        <div>
          <img [src]="user.icon.url" class="w-24 h-24 rounded-lg" />
        </div>
        <div class="flex flex-col items-start py-2">
          <span class="font-bold text-xl">{{ user.name }}</span>
          <span class="text-sm text-gray-600">@{{ user.id }}@{{ hostname }}</span>
        </div>
        <div class="py-2" [innerHTML]="user.description"></div>
      </div>
    </ng-container>
  `,
  styles: [],
})
export class UserComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly http = inject(HttpClient);
  private readonly state = new RxState<{ user: LocalUser | null }>();
  readonly state$ = this.state.select().pipe(stateful());
  readonly username$ = this.route.params.pipe(map((params) => params['username']));

  readonly hostname = window.location.hostname;

  ngOnInit() {
    this.state.connect(
      'user',
      this.username$.pipe(
        switchMap((username) => this.http.get<{ user: LocalUser }>(`/api/users/show/${username}`)),
        map((resp) => resp.user),
      ),
    );
  }
}
