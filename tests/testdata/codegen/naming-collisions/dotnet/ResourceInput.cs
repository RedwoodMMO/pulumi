// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Example
{
    [ExampleResourceType("example::ResourceInput")]
    public partial class ResourceInput : global::Pulumi.CustomResource
    {
        [Output("bar")]
        public Output<string?> Bar { get; private set; } = null!;


        /// <summary>
        /// Create a ResourceInput resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public ResourceInput(string name, ResourceInputArgs? args = null, CustomResourceOptions? options = null)
            : base("example::ResourceInput", name, args ?? new ResourceInputArgs(), MakeResourceOptions(options, ""))
        {
        }

        private ResourceInput(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("example::ResourceInput", name, null, MakeResourceOptions(options, id))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
        /// <summary>
        /// Get an existing ResourceInput resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static ResourceInput Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new ResourceInput(name, id, options);
        }
    }

    public sealed class ResourceInputArgs : global::Pulumi.ResourceArgs
    {
        public ResourceInputArgs()
        {
        }
        public static new ResourceInputArgs Empty => new ResourceInputArgs();
    }
}