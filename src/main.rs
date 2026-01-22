//! An example showing a very basic implementation.

use avian2d::prelude::*;
use bevy::prelude::*;
use bevy_simple_text_input::{
    TextInput, TextInputPlugin, TextInputSubmitMessage, TextInputSystem, TextInputTextColor,
    TextInputTextFont,
};
use lat::{Guest, Lat};

#[derive(Component)]
struct SandboxEntity;

const BORDER_COLOR_ACTIVE: Color = Color::srgb(0.75, 0.52, 0.99);
const TEXT_COLOR: Color = Color::srgb(0.9, 0.9, 0.9);
const BACKGROUND_COLOR: Color = Color::srgb(0.15, 0.15, 0.15);

fn main() {
    App::new()
        .add_plugins((
            DefaultPlugins.set(WindowPlugin {
                primary_window: Some(Window {
                    resolution: bevy::window::WindowResolution::new(1280, 720),
                    resizable: false,
                    ..default()
                }),
                ..default()
            }),
            PhysicsPlugins::default().with_length_unit(1.0),
        ))
        .add_plugins(TextInputPlugin)
        .insert_resource(Gravity(Vec2::NEG_Y * 981.0))
        .init_resource::<Messages<ResetMessage>>()
        .add_systems(Startup, setup)
        .add_systems(Update, (listener.after(TextInputSystem), update_boundaries))
        .add_systems(Update, reset)
        .run();
}

fn setup(
    mut commands: Commands,
    mut meshes: ResMut<Assets<Mesh>>,
    mut materials: ResMut<Assets<ColorMaterial>>,
) {
    commands.spawn((Camera2d::default(),));

    commands
        .spawn(Node {
            width: Val::Percent(100.0),
            height: Val::Percent(100.0),
            justify_content: JustifyContent::SpaceBetween,
            align_items: AlignItems::Start,
            padding: UiRect::all(Val::Px(5.0)),
            ..default()
        })
        .with_children(|parent| {
            parent.spawn((
                Node {
                    border: UiRect::all(Val::Px(5.0)),
                    padding: UiRect::all(Val::Px(5.0)),
                    ..default()
                },
                TextInput,
                TextInputTextFont(TextFont {
                    font_size: 34.,
                    ..default()
                }),
                TextInputTextColor(TextColor(TEXT_COLOR)),
                BorderColor::all(BORDER_COLOR_ACTIVE),
                BackgroundColor(BACKGROUND_COLOR),
            ));
        });

    // 1. Boundary (Walls/Floor)
    // Define screen dimensions (example: 1280x720)
    let w = 1280.0;
    let h = 720.0;
    let vertices = vec![
        Vec2::new(-w, -h),
        Vec2::new(w, -h),
        Vec2::new(w, h),
        Vec2::new(-w, h),
        Vec2::new(-w, -h),
    ];

    commands.spawn((
        RigidBody::Static,
        Collider::polyline(vertices, None),
        Boundary,
        SandboxEntity,
    ));

    // 2. Falling Ball
    let radius = 20.0;
    commands.spawn((
        Mesh2d(meshes.add(Circle::new(radius))),
        MeshMaterial2d(materials.add(ColorMaterial::from_color(Color::srgb(1.0, 0.0, 0.0)))),
        Transform::from_xyz(0.0, 0.0, 0.0),
        RigidBody::Dynamic,
        Collider::circle(radius),
        Mass(1.0),
        Restitution::new(0.6), // Bounce
        Name::new("Ball"),
        SandboxEntity,
    ));
}

#[derive(Component)]
struct Boundary;

fn update_boundaries(windows: Query<&Window>, mut query: Query<&mut Collider, With<Boundary>>) {
    let window = windows.single().unwrap();
    let w = window.width() / 2.0;
    let h = window.height() / 2.0;

    if let Ok(mut collider) = query.single_mut() {
        // Recreate polyline based on new window dimensions
        *collider = Collider::polyline(
            vec![
                Vec2::new(-w, -h),
                Vec2::new(w, -h),
                Vec2::new(w, h),
                Vec2::new(-w, h),
                Vec2::new(-w, -h),
            ],
            None,
        );
    }
}

fn listener(
    mut events: MessageReader<TextInputSubmitMessage>,
    mut resetter: MessageWriter<ResetMessage>,
) {
    for event in events.read() {
        if let Ok(result) = Lat::parse(event.value.clone())
            && result.resetting
        {
            resetter.write(ResetMessage);
        }
    }
}

#[derive(Message)]
struct ResetMessage;

fn reset(
    mut commands: Commands,
    mut meshes: ResMut<Assets<Mesh>>,
    mut materials: ResMut<Assets<ColorMaterial>>,
    mut resetter: MessageReader<ResetMessage>,
    sandbox_entities: Query<Entity, With<SandboxEntity>>,
) {
    for _reset_msg in resetter.read() {
        for entity in sandbox_entities.into_iter() {
            commands.entity(entity).despawn();
        }

        // 1. Boundary (Walls/Floor)
        // Define screen dimensions (example: 1280x720)
        let w = 1280.0;
        let h = 720.0;
        let vertices = vec![
            Vec2::new(-w, -h),
            Vec2::new(w, -h),
            Vec2::new(w, h),
            Vec2::new(-w, h),
            Vec2::new(-w, -h),
        ];

        commands.spawn((
            RigidBody::Static,
            Collider::polyline(vertices, None),
            SandboxEntity,
            Boundary,
        ));

        // 2. Falling Ball
        let radius = 20.0;
        commands.spawn((
            Mesh2d(meshes.add(Circle::new(radius))),
            MeshMaterial2d(materials.add(ColorMaterial::from_color(Color::srgb(1.0, 0.0, 0.0)))),
            Transform::from_xyz(0.0, 0.0, 0.0),
            RigidBody::Dynamic,
            Collider::circle(radius),
            Mass(1.0),
            Restitution::new(0.6), // Bounce
            Name::new("Ball"),
            SandboxEntity,
        ));
    }
}
