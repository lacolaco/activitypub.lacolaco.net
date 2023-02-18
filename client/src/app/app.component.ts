import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { Auth, authState, GoogleAuthProvider, signInWithPopup, signOut, User } from '@angular/fire/auth';
import { RouterLink, RouterOutlet } from '@angular/router';
import { RxState, stateful } from '@rx-angular/state';
import { tap } from 'rxjs';
import { AppStrokedButton } from './shared/ui/button';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterOutlet, AppStrokedButton],
  template: `
    <ng-container *ngIf="state$ | async as state">
      <header class="p-4 shadow-gray-500 shadow-sm z-10">
        <div class="container flex flex-row items-center justify-between">
          <h1 class="font-bold">activitypub.lacolaco.net</h1>
          <div *ngIf="!state.user">
            <button app-stroked-button (click)="signIn()" class="text-sm">Sign in</button>
          </div>
          <div *ngIf="state.user">
            <button app-stroked-button (click)="signOut()" class="text-sm">Sign out</button>
          </div>
        </div>
      </header>
      <main class="flex-auto container py-4 flex flex-col gap-y-2">
        <div *ngIf="state.user" class="flex flex-col">
          <a routerLink="/search" app-stroked-button>Control Room</a>
        </div>
        <router-outlet></router-outlet>
      </main>
    </ng-container>
  `,
  host: { class: 'flex flex-col w-full h-full bg-white font-sans' },
})
export class AppComponent {
  private readonly auth = inject(Auth);

  private readonly state = new RxState<{ user: User | null }>();
  readonly state$ = this.state.select().pipe(stateful());

  ngOnInit() {
    this.state.set({ user: null });
    this.state.connect(
      'user',
      authState(this.auth).pipe(
        tap((user) => {
          console.log(user);
        }),
      ),
    );
  }

  async signIn() {
    try {
      await signInWithPopup(this.auth, new GoogleAuthProvider());
    } catch (e) {
      console.error(e);
    }
  }

  async signOut() {
    try {
      await signOut(this.auth);
    } catch (e) {
      console.error(e);
    }
  }
}
