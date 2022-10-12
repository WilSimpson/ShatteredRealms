import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {HomeComponent} from './pages/home/home.component';
import {LoginComponent} from './pages/login/login.component';
import {RegisterComponent} from "./pages/register/register.component";
import {UserProfileComponent} from "./pages/user-profile/user-profile.component";
import {AdminDashboardComponent} from "./pages/admin-dashboard/admin-dashboard.component";
import {AuthGuard} from "./_helpers/auth.guard";
import {UserEditComponent} from "./pages/user-edit/user-edit.component";
import {ViewRoleComponent} from "./pages/roles/view-role/view-role.component";
import {NewRoleComponent} from "./pages/roles/new-role/new-role.component";
import {EditRoleComponent} from "./pages/roles/edit-role/edit-role.component";

const routes: Routes = [
  {path: '', component: HomeComponent},
  {path: 'login', component: LoginComponent},
  {path: 'register', component: RegisterComponent},
  {path: 'users/:user', component: UserProfileComponent},
  {path: 'users/:user/edit', component: UserEditComponent},
  {path: 'roles/new', component: NewRoleComponent},
  {path: 'roles/:role', component: ViewRoleComponent},
  {path: 'roles/:role/edit', component: EditRoleComponent},
  {
    path: 'admin',
    component: AdminDashboardComponent,
    data: {roles: ["SUPER ADMIN", "ADMIN"]},
    canActivate: [AuthGuard],
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
