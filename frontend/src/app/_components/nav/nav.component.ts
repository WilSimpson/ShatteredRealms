import { Component, OnInit } from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {AuthenticationService} from "../../_services/authentication.service";
import {fromEvent} from "rxjs";
import {debounceTime, map, startWith} from "rxjs/operators";

@Component({
  selector: 'app-nav',
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss']
})
export class NavComponent implements OnInit {
  links = [
    {title: 'Home', link: '/'},
    {title: 'Login', link: '/login'},
  ];

  isScreenSmall$: any;

  constructor(protected route: ActivatedRoute,
              protected router: Router,
              public authService: AuthenticationService) { }

  ngOnInit(): void {
    // Checks if screen size is less than 1024 pixels
    const checkScreenSize = () => document.body.offsetWidth < 992;

    // Create observable from window resize event throttled so only fires every 500ms
    const screenSizeChanged$ = fromEvent(window, 'resize').pipe(debounceTime(200), map(checkScreenSize));

    // Start off with the initial value use the isScreenSmall$ | async in the
    // view to get both the original value and the new value after resize.
    this.isScreenSmall$ = screenSizeChanged$.pipe(startWith(checkScreenSize()))
  }

  isUserSignedIn(): boolean {
    return this.authService.currentUserValue != undefined;
  }

  signOutUser(): void {
    this.authService.logout();
    this.router.navigate(['/login']).then();
  }
}
