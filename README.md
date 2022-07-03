# Intro
THIS IS A EXPERIMENT - DO NOT USE!

This is a NMEA0183 parser/printer.

The spec is based on https://gpsd.gitlab.io/gpsd/NMEA.html


## Code generation

spec.yaml --> transform --> render --> compile
                 ^            ^
                 |            |
              impl specific  templates
              data

Transformation prepares the data for easy rendering.
1. Map spec types to implementation types, for example Float maps to float64.
2. Combine fields in single value


### Related fields

Multiple related fields can be combined into one entity, for example a number field and a string field into a Distance entity.

In the spec related fields start with the same name but have different '.Xyz' suffixes.
```yaml
fields:
- name: WaypointRadius
  type: Radius
- name: WaypointRadius.Unit
```


#### Parsing

The parser methods is passed one or more field indices. 
The first is of the 'entity', followed by sub-fields in order of appearanace (left to right)


#### Printing 


