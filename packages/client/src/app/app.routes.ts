import { Routes, UrlSegment } from '@angular/router';
import { UserComponent } from './user/user.component';
import { UsersComponent } from './users/users.component';

export const routes: Routes = [
  {
    path: 'users',
    component: UsersComponent,
  },
  {
    path: 'users/:hostname/:username',
    component: UserComponent,
  },
  {
    path: '**',
    redirectTo: '/users',
  },
];
