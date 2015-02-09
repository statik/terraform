package aws

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/goamz/ec2"
	"github.com/mitchellh/goamz/autoscaling"
)

// tagsSchema returns the schema to use for tags.
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

// setTags is a helper to set the tags for a resource. It expects the
// tags field to be named "tags"
func setTags(conn *ec2.EC2, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		// Set tags
		if len(remove) > 0 {
			log.Printf("[DEBUG] Removing tags: %#v", remove)
			if _, err := conn.DeleteTags([]string{d.Id()}, remove); err != nil {
				return err
			}
		}
		if len(create) > 0 {
			log.Printf("[DEBUG] Creating tags: %#v", create)
			if _, err := conn.CreateTags([]string{d.Id()}, create); err != nil {
				return err
			}
		}
	}

	return nil
}

// diffTags takes our tags locally and the ones remotely and returns
// the set of tags that must be created, and the set of tags that must
// be destroyed.
func diffTags(oldTags, newTags []ec2.Tag) ([]ec2.Tag, []ec2.Tag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []ec2.Tag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return tagsFromMap(create), remove
}

// tagsFromMap returns the tags for the given map of data.
func tagsFromMap(m map[string]interface{}) []ec2.Tag {
	result := make([]ec2.Tag, 0, len(m))
	for k, v := range m {
		result = append(result, ec2.Tag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

// tagsToMap turns the list of tags into a map.
func tagsToMap(ts []ec2.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range ts {
		result[t.Key] = t.Value
	}

	return result
}

func autoscalingTagsToMap(ts []autoscaling.Tag) map[string]struct{string;bool} {
  result := make(map[string]struct{string;bool})
  for _, t := range ts {
    result[t.Key] = struct{string;bool}{t.Value, t.PropagateAtLaunch}
  }
  return result
}
