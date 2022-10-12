import { Component, OnInit } from '@angular/core';
import {AuthenticationService} from "../../_services/authentication.service";
import {User} from "../../models/user.model";

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent implements OnInit {

  constructor(protected authService: AuthenticationService) { }

  ngOnInit(): void {
  }

  public get user(): User | null {
    return this.authService.currentUserValue;
  }

  isUserSignedIn(): boolean {
    return this.user != null;
  }

  isAdmin(): boolean {
    return this.authService.hasAnyRole(["SUPER ADMIN", "ADMIN"])
  }
}
