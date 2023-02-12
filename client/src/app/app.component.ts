import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { Auth, authState, GoogleAuthProvider, signInWithPopup, signOut, User } from '@angular/fire/auth';
import { RouterOutlet } from '@angular/router';
import { RxState, stateful } from '@rx-angular/state';
import { tap } from 'rxjs';
import { AppStrokedButton } from './shared/ui/button';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, AppStrokedButton],
  template: `
    <header class="p-4 shadow-gray-500 shadow-sm z-10">
      <div class="container">
        <span style="display: block">{{ title }} app is running!</span>
      </div>
    </header>
    <main class="flex-auto container py-4">
      <ng-container *ngIf="state$ | async as state">
        <div *ngIf="!state.user" class="w-full py-2">
          <button app-stroked-button (click)="signIn()">Sign in</button>
        </div>
        <div *ngIf="state.user" class="w-full py-2 flex flex-row items-center gap-x-2">
          <button app-stroked-button (click)="signOut()">Sign out</button>
          <div>Current User: {{ state.user.displayName }}</div>
        </div>
        <router-outlet> </router-outlet>
      </ng-container>
    </main>
  `,
  styles: [
    `
      :host {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
      }
    `,
  ],
})
export class AppComponent {
  title = 'activitypub.lacolaco.net';

  private readonly auth = inject(Auth);

  private readonly state = new RxState<{ user: User | null }>();
  readonly state$ = this.state.select().pipe(stateful());

  ngOnInit() {
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
