import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NavComponent } from './nav.component';
import {ActivatedRoute, Router} from "@angular/router";
import {RouterTestingModule} from "@angular/router/testing";
import {MdbCollapseModule} from "mdb-angular-ui-kit/collapse";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AuthenticationService} from "../../_services/authentication.service";

describe('NavComponent', () => {
  let component: NavComponent;
  let fixture: ComponentFixture<NavComponent>;
  let route: ActivatedRoute;
  let authService: AuthenticationService;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ NavComponent ],
      imports: [RouterTestingModule, MdbCollapseModule, HttpClientTestingModule],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(NavComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    route = TestBed.inject(ActivatedRoute);
    authService = TestBed.inject(AuthenticationService);
    router = TestBed.inject(Router);
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should know if a user is logged in', () => {
    spyOnProperty(authService, 'currentUserValue').and.returnValue({});
    expect(component.isUserSignedIn()).toBeTrue();
  });


  it('should be able to signout a user', () => {
    spyOn(authService, 'logout');
    spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));
    component.signOutUser();

    expect(authService.logout).toHaveBeenCalledTimes(1);
    expect(router.navigate).toHaveBeenCalledTimes(1);
  });
});
