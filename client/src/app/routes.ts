import { Routes } from '@angular/router';
import { FooComponent } from './foo/foo.component';

export const routes: Routes = [{ path: 'users/:username', component: FooComponent }];
